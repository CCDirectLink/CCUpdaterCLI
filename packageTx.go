package ccmodupdater

import (
	"fmt"
	"sort"
	"github.com/Masterminds/semver"
)

// It slices! It dices! It lops off half of your hand when you're looking at the telly,
//  because that's what happens when you mishandle industrial equipment!
// It's the Package Transaction Engine!

// A PackageTXOperation describes something to perform on a package.
type PackageTXOperation int
const (
	// PackageTransactionOperationInstall updates the package to the latest version.
	PackageTXOperationInstall PackageTXOperation = iota
	// PackageTransactionOperationRemove removes the package from the system.
	PackageTXOperationRemove
)

// A PackageTX is a proposed set of operations to perform on a GameInstance.
// It should be treated immutably, using clone-and-edit where required.
type PackageTX map[string]PackageTXOperation

// Clone clones the PackageTX.
func (tx PackageTX) Clone() PackageTX {
	base := map[string]PackageTXOperation{}
	for k, v := range tx {
		base[k] = v
	}
	return base
}

// Append creates a PackageTX from two PackageTXes, showing the result of them having been concatenated.
func (tx PackageTX) Append(b PackageTX) PackageTX {
	base := tx.Clone()
	for k, v := range b {
		base[k] = v
	}
	return base
}

// Equals checks for equality with another PackageTX.
func (tx PackageTX) Equals(other PackageTX) bool {
	for k, v := range tx {
		otherOp, present := other[k]
		if (!present) || (otherOp != v)  {
			return false
		}
	}
	for k, v := range other {
		txOp, present := tx[k]
		if (!present) || (txOp != v)  {
			return false
		}
	}
	return true
}

// PackageTXSet represents a set of PackageTXes. Refer to it via *PackageTXSet.
type PackageTXSet []PackageTX
// Contains returns true if the PackageTX is contained in this set.
func (set *PackageTXSet) Contains(ptx PackageTX) bool {
	for _, v := range *set {
		if v.Equals(ptx) {
			return true
		}
	}
	return false
}
// Put puts a PackageTX into the set. Returns false if it was already there.
func (set *PackageTXSet) Put(ptx PackageTX) bool {
	if !set.Contains(ptx) {
		*set = append(*set, ptx)
		return true
	}
	return false
}

// Len implements sort.Interface.Len
func (set *PackageTXSet) Len() int {
	return len(*set)
}

// Swap implements sort.Interface.Swap
func (set *PackageTXSet) Swap(a int, b int) {
	av := (*set)[a]
	(*set)[a] = (*set)[b]
	(*set)[b] = av
}

// Less implements sort.Interface.Less
func (set *PackageTXSet) Less(a int, b int) bool {
	av := (*set)[a]
	bv := (*set)[b]
	return len(av) < len(bv)
}

// PackageTXProblem represents a problem with the PackageTX and a number of potential solutions.
type PackageTXProblem struct {
	Text string
	// Solutions. Each solution is a PackageTX to be appended to the original.
	Solutions []PackageTX
}

// A PackageTXContext is a view of the current local and remote packages from before a package transaction.
type PackageTXContext struct {
	LocalPackages map[string]LocalPackage
	// Leave non-nil but blank if offline.
	RemotePackages map[string]RemotePackage
}
func (ctx PackageTXContext) AllPackageIDs() []string {
	set := []string{}
	for k := range ctx.LocalPackages {
		set = append(set, k)
	}
	for k := range ctx.RemotePackages {
		if ctx.LocalPackages[k] == nil {
			set = append(set, k)
		}
	}
	return set
}

func (ctx PackageTXContext) PackageAfter(tx PackageTX, pkg string) Package {
	op, present := tx[pkg]
	if present {
		if op == PackageTXOperationInstall {
			return ctx.RemotePackages[pkg]
		} else if op == PackageTXOperationRemove {
			return nil
		} else {
			panic("no idea what this operation does to pkg " + pkg + " in PackageAfter")
		}
	}
	return ctx.LocalPackages[pkg]
}

// Verify checks if the Package Transaction is invalid, and if so, why.
func (ctx PackageTXContext) Verify(tx PackageTX) []PackageTXProblem {
	// PackageTXContexts: CCMIAs In Disguise
	problems := []PackageTXProblem{}
	for pkg, op := range tx {
		localPkg := ctx.LocalPackages[pkg]
		remotePkg := ctx.RemotePackages[pkg]
		if op == PackageTXOperationRemove {
			if localPkg != nil {
				// Find packages dependent on it that are not also being removed
				for _, depID := range ctx.AllPackageIDs() {
					dep := ctx.PackageAfter(tx, depID)
					if dep != nil {
						depVersionDep := dep.Metadata().Dependencies()[pkg]
						if depVersionDep != nil {
							problems = append(problems, PackageTXProblem{
								Text: localPkg.Metadata().FullReferent() + " cannot be removed because it is required by " + dep.Metadata().FullReferent() + ".",
								Solutions: []PackageTX{},
							})
						}
					}
				}
			} else {
				problems = append(problems, PackageTXProblem{
					Text: pkg + " cannot be removed because it is not installed.",
					Solutions: []PackageTX{},
				})
			}
		} else if op == PackageTXOperationInstall {
			if remotePkg != nil {
				for depID, depC := range remotePkg.Metadata().Dependencies() {
					dep := ctx.PackageAfter(tx, depID)
					if dep == nil {
						problems = append(problems, PackageTXProblem{
							Text: depID + " is required by " + pkg + " but is not installed.",
							Solutions: []PackageTX{
								PackageTX{
									depID: PackageTXOperationInstall,
								},
							},
						})
					} else {
						if !depC.Check(dep.Metadata().Version()) {
							problems = append(problems, PackageTXProblem{
								Text: fmt.Sprintf("%s is version %s, which does not meet requirement from %s: %s", dep.Metadata().FullReferent(), dep.Metadata().Version().String(), remotePkg.Metadata().FullReferent(), depC),
								Solutions: []PackageTX{
									PackageTX{
										depID: PackageTXOperationInstall,
									},
								},
							})
						}
					}
				}
			} else {
				problems = append(problems, PackageTXProblem{
					Text: pkg + " cannot be retrieved because there is no remote package.",
					Solutions: []PackageTX{},
				})
			}
		} else {
			problems = append(problems, PackageTXProblem{
				Text: pkg + " cannot be operated on because the solver did not understand the command.",
				Solutions: []PackageTX{},
			})
		}
	}
	return problems
}

// Solve returns a set of solved versions of the package transactions.
func (ctx PackageTXContext) Solve(tx PackageTX) ([]PackageTX, error) {
	// solutions is the set of all completed PackageTXes.
	solutions := &PackageTXSet{}
	// explored is the set of all PackageTXes explored.
	explored := &PackageTXSet{}
	// toExplore is the set of PackageTXes to explore this round.
	toExplore := &PackageTXSet{tx}
	
	for len(*toExplore) > 0 {
		// Round!
		lastToExplore := toExplore
		toExplore = &PackageTXSet{}
		for _, etx := range *lastToExplore {
			if !explored.Put(etx) {
				// Already explored!
				continue
			}
			problems := ctx.Verify(etx)
			if len(problems) == 0 {
				solutions.Put(etx)
				continue
			}
			for _, problem := range problems {
				for _, solution := range problem.Solutions {
					// Append the current state with the solution, and then the original TX for goal preservation.
					toExplore.Put(etx.Append(solution).Append(tx))
				}
			}
		}
	}
	if len(*solutions) > 0 {
		sort.Sort(solutions)
		return *solutions, nil
	}
	
	// Ok, so we failed; *why* did we fail?
	originalProblems := ctx.Verify(tx)
	originalProblemDetail := ""
	for _, v := range originalProblems {
		originalProblemDetail += "\n" + v.Text
	}
	return nil, fmt.Errorf("Unsolvable packaging situation%s", originalProblemDetail)
}

func (ctx PackageTXContext) findUnresolved(deps map[string]*semver.Constraints) int {
	counter := 0
	for k, v := range deps {
		ctxP := ctx.LocalPackages[k]
		if ctxP == nil {
			counter++
			continue
		}
		if !v.Check(ctxP.Metadata().Version()) {
			counter++
		}
	}
	return counter
}

// Perform actually performs a PackageTX (invalidating the PackageTXContext LocalPackages as a result) in dependency order.
func (ctx PackageTXContext) Perform(gi *GameInstance, tx PackageTX, stats func (pkg string, pre bool, rm bool, add bool)) error {
	// Make a copy that we can modify safely.
	tx = tx.Clone()
	for true {
		hasBest := false
		best := ""
		for k, op := range tx {
			// Rules:
			// Firstly, execute installs from lowest amounts of unresolved dependencies upwards.
			// Secondly, perform removals.
			if !hasBest {
				best = k
				hasBest = true
			} else if op == PackageTXOperationInstall {
				if tx[best] != PackageTXOperationInstall {
					// Install overrides any other operation
					best = k
				} else {
					// Between two Install operations, install the one with the least unresolved deps
					remoteBest := ctx.RemotePackages[best]
					if remoteBest == nil {
						return fmt.Errorf("No remote package: %s", remoteBest)
					}
					remoteK := ctx.RemotePackages[k]
					if remoteK == nil {
						return fmt.Errorf("No remote package: %s", remoteK)
					}
					unresolvedBest := ctx.findUnresolved(remoteBest.Metadata().Dependencies())
					unresolvedK := ctx.findUnresolved(remoteK.Metadata().Dependencies())
					if unresolvedK < unresolvedBest {
						best = k
					}
				}
			}
		}
		if !hasBest {
			break
		}
		currentOp := tx[best]
		// Be warned! This is what prevents an infinite loop.
		delete(tx, best)
		
		if currentOp == PackageTXOperationRemove {
			pkg := ctx.LocalPackages[best]
			if pkg == nil {
				return fmt.Errorf("%s to be removed but is non-existent", best)
			}
			stats(best, true, true, false)
			err := pkg.Remove()
			if err != nil {
				return fmt.Errorf("During removal of %s: %s", pkg.Metadata().FullReferent(), err.Error())
			}
			stats(best, false, true, false)
		} else if currentOp == PackageTXOperationInstall {
			pkg := ctx.LocalPackages[best]
			stats(best, true, pkg != nil, true)
			pkgRemote := ctx.RemotePackages[best]
			if pkgRemote == nil {
				// Run this check BEFORE the package is potentially removed
				return fmt.Errorf("No remote package: %s", best)
			}
			if pkg != nil {
				// Remove in preparation for upgrade
				err := pkg.Remove()
				if err != nil {
					return fmt.Errorf("%s could not be removed for upgrade: %s", pkg.Metadata().FullReferent(), err.Error())
				}
			}
			err := pkgRemote.Install(gi)
			if err != nil {
				stats(best, false, pkg != nil, false)
				return fmt.Errorf("%s could not be installed: %s", pkgRemote.Metadata().FullReferent(), err.Error())
			}
			stats(best, false, pkg != nil, true)
		} else {
			return fmt.Errorf("unable to handle operation type for %s", best)
		}
		// Update ctx.LocalPackages (as this isn't a pointer, won't help caller)
		ctx.LocalPackages = gi.Packages()
	}
	return nil
}
