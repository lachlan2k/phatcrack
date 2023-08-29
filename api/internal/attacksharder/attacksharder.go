package attacksharder

import (
	"fmt"

	"github.com/lachlan2k/phatcrack/api/internal/db"
	"github.com/lachlan2k/phatcrack/api/internal/hashcathelpers"
	"github.com/lachlan2k/phatcrack/common/pkg/hashcattypes"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Ref: https://hashcat.net/wiki/doku.php?id=mask_attack
var builtinCharset = map[byte][]byte{
	'l': []byte("abcdefghijklmnopqrstuvwxyz"),
	'u': []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),

	'd': []byte("0123456789"),
	'h': []byte("0123456789abcdef"),
	'H': []byte("0123456789ABCDEF"),

	's': []byte(" !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"),

	'a': []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"),
	'b': {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255},
}

func splitMask(inputMask string, numChunks int) (outputMask string, outputCharsets [][]byte, err error) {
	firstVariableIndex := -1
	for i := range inputMask {
		if inputMask[i] != '?' {
			continue
		}

		if i == len(inputMask)-1 {
			return "", nil, fmt.Errorf("'%q' is not a valid mask (unexpected ? at end)", inputMask)
		}

		firstVariableIndex = i
		break
	}

	// While a mask with no variables in it is dumb, its not *technically* wrong
	// But I'm going to throw a fit regardless
	if firstVariableIndex == -1 {
		return "", nil, fmt.Errorf("'%q' cannot be chunked as it contains no variables", inputMask)
	}

	maskChar := inputMask[firstVariableIndex+1]
	builtinCharset, ok := builtinCharset[maskChar]
	if !ok {
		return "", nil, fmt.Errorf("the first variable in the mask '%q' was not recognized", inputMask)
	}

	// Use ?4 as the "magic" custom charset dedicated to sharding
	outputMask = inputMask[:firstVariableIndex+1] + "4" + inputMask[firstVariableIndex+2:]
	outputCharsets = [][]byte{}

	for i := 0; i < numChunks; i++ {
		start := (i * len(builtinCharset) / numChunks)
		end := ((i + 1) * len(builtinCharset)) / numChunks

		outputCharsets = append(outputCharsets, builtinCharset[start:end])
	}

	return outputMask, outputCharsets, nil
}

func createSingleJobfromAttack(attack *db.Attack) (*db.Job, *db.Hashlist, error) {
	hashlist, err := db.GetHashlistWithHashes(attack.HashlistID.String())
	if err != nil {
		return nil, nil, err
	}

	targetHashes := []string{}
	for _, hash := range hashlist.Hashes {
		if !hash.IsCracked {
			targetHashes = append(targetHashes, hash.NormalizedHash)
		}
	}

	dbJob, err := db.CreateJob(&db.Job{
		HashlistVersion: hashlist.Version,
		AttackID:        &attack.ID,
		HashcatParams:   attack.HashcatParams,
		TargetHashes:    targetHashes,
		HashType:        hashlist.HashType,
	})

	if err != nil {
		return nil, nil, err
	}

	return dbJob, hashlist, err
}

func shardMaskAttack(attack *db.Attack, numJobs int) ([]*db.Job, *db.Hashlist, error) {
	params := attack.HashcatParams.Data
	if len(params.MaskCustomCharsets) >= 4 {
		return nil, nil, fmt.Errorf("received %d custom character sets, maximum is 3", len(params.MaskCustomCharsets))
	}

	outputMask, shardedCharsets, err := splitMask(params.Mask, numJobs)
	if err != nil {
		return nil, nil, err
	}

	hashlist, err := db.GetHashlistWithHashes(attack.HashlistID.String())
	if err != nil {
		return nil, nil, err
	}

	targetHashes := []string{}
	for _, hash := range hashlist.Hashes {
		if !hash.IsCracked {
			targetHashes = append(targetHashes, hash.NormalizedHash)
		}
	}

	params.Mask = outputMask
	jobs := []*db.Job{}

	err = db.GetInstance().Transaction(func(tx *gorm.DB) error {
		for i := 0; i < numJobs; i++ {
			params.MaskShardedCharset = string(shardedCharsets[i])

			dbJob, err := db.CreateJobTx(&db.Job{
				HashlistVersion: hashlist.Version,
				AttackID:        &attack.ID,
				HashcatParams:   datatypes.NewJSONType(params),
				TargetHashes:    targetHashes,
				HashType:        hashlist.HashType,
			}, tx)

			jobs = append(jobs, dbJob)

			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return jobs, hashlist, nil
}

func shardAttackByKeyspace(attack *db.Attack, numJobs int) ([]*db.Job, *db.Hashlist, error) {
	params := attack.HashcatParams.Data
	hashlist, err := db.GetHashlistWithHashes(attack.HashlistID.String())
	if err != nil {
		return nil, nil, err
	}

	targetHashes := []string{}
	for _, hash := range hashlist.Hashes {
		if !hash.IsCracked {
			targetHashes = append(targetHashes, hash.NormalizedHash)
		}
	}

	keyspace, err := hashcathelpers.CalculateKeyspace(params)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't calculate keyspace for sharding: %w", err)
	}

	limitPerJob := keyspace / int64(numJobs)
	jobs := []*db.Job{}

	err = db.GetInstance().Transaction(func(tx *gorm.DB) error {
		for i := 0; i < numJobs; i++ {
			params.Skip = limitPerJob * int64(i)

			if i == numJobs-1 {
				params.Limit = 0
			} else {
				params.Limit = limitPerJob
			}

			dbJob, err := db.CreateJobTx(&db.Job{
				HashlistVersion: hashlist.Version,
				AttackID:        &attack.ID,
				HashcatParams:   datatypes.NewJSONType(params),
				TargetHashes:    targetHashes,
				HashType:        hashlist.HashType,
			}, tx)

			jobs = append(jobs, dbJob)

			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return jobs, hashlist, nil
}

func MakeJobs(attack *db.Attack, maxNumJobs int) ([]*db.Job, *db.Hashlist, error) {
	if !attack.IsDistributed || maxNumJobs <= 1 {
		job, h, err := createSingleJobfromAttack(attack)
		if err != nil {
			return nil, nil, err
		}
		jobs := []*db.Job{job}
		return jobs, h, err
	}

	switch attack.HashcatParams.Data.AttackMode {
	case hashcattypes.AttackModeDictionary, hashcattypes.AttackModeCombinator:
		return shardAttackByKeyspace(attack, maxNumJobs)

	case hashcattypes.AttackModeMask, hashcattypes.AttackModeHybridDM, hashcattypes.AttackModeHybridMD:
		return shardMaskAttack(attack, maxNumJobs)

	default:
		return nil, nil, fmt.Errorf("unrecognized hash type: %d", attack.HashcatParams.Data.HashType)
	}
}
