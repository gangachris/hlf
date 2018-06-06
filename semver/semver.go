package semver

import (
	"errors"
	"strconv"
	"strings"
)

type semanticVersion struct {
	major, minor, patch int
}

var errInvalidSemver = errors.New("invalid semver")

func getSemverFromString(semver string) (semanticVersion, error) {
	smv := semanticVersion{}

	splitSemver := strings.Split(semver, ".")
	if len(splitSemver) == 0 {
		return smv, errInvalidSemver
	}

	if len(splitSemver) == 3 {
		return setSemverFromSlice(splitSemver)
	}

	var ss []string
	if len(splitSemver) == 2 {
		ss = append(ss, splitSemver...)
		ss = append(ss, "0")
		return setSemverFromSlice(ss)
	}

	if len(splitSemver) == 1 {
		ss = append(ss, splitSemver...)
		ss = append(ss, "0", "0")
		return setSemverFromSlice(ss)
	}

	return smv, errInvalidSemver
}

func setSemverFromSlice(semver []string) (semanticVersion, error) {
	smv := semanticVersion{}
	major, err := strconv.Atoi(semver[0])
	if err != nil {
		return smv, err
	}
	smv.major = major

	minor, err := strconv.Atoi(semver[1])
	if err != nil {
		return smv, err
	}
	smv.minor = minor

	patch, err := strconv.Atoi(semver[2])
	if err != nil {
		return smv, err
	}
	smv.patch = patch

	return smv, nil
}

// CorrectVersion compares the two semantic versions provided and returns
// true if current version is greater than or equal to the minimum version
func CorrectVersion(minimum, current string) (bool, error) {
	minSemver, err := getSemverFromString(minimum)
	if err != nil {
		return false, err
	}

	currSemver, err := getSemverFromString(current)
	if err != nil {
		return false, err
	}

	// compare
	if minSemver.major != currSemver.major {
		return minSemver.major < currSemver.major, nil
	}

	if minSemver.minor != currSemver.minor {
		return minSemver.minor < currSemver.minor, nil
	}

	if minSemver.patch != currSemver.patch {
		return minSemver.patch < currSemver.patch, nil
	}

	return minSemver.patch == currSemver.patch, nil
}
