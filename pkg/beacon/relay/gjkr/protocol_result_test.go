package gjkr

import (
	"reflect"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
)

func TestGenerateResult(t *testing.T) {
	threshold := 3
	groupSize := 8

	var tests = map[string]struct {
		disqualifiedMemberIDs []group.MemberIndex
		inactiveMemberIDs     []group.MemberIndex
		expectedResult        func(groupPublicKey *bn256.G2) *Result
	}{
		"no disqualified or inactive members - success": {
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: groupPublicKey,
					Group:          initalizeGroup(groupSize, []group.MemberIndex{}, []group.MemberIndex{}),
				}
			},
		},
		"one disqualified member - success": {
			disqualifiedMemberIDs: []group.MemberIndex{2},
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: groupPublicKey,
					Group:          initalizeGroup(groupSize, []group.MemberIndex{2}, []group.MemberIndex{}),
				}
			},
		},
		"two inactive members - success": {
			inactiveMemberIDs: []group.MemberIndex{3, 7},
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: groupPublicKey,
					Group:          initalizeGroup(groupSize, []group.MemberIndex{}, []group.MemberIndex{3, 7}),
				}
			},
		},
		"more than half of threshold disqualified and inactive members - failure": {
			disqualifiedMemberIDs: []group.MemberIndex{2},
			inactiveMemberIDs:     []group.MemberIndex{3, 7},
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: nil,
					Group:          initalizeGroup(groupSize, []group.MemberIndex{2}, []group.MemberIndex{3, 7}),
				}
			},
		},
		"more than half of threshold inactive members - failure": {
			inactiveMemberIDs: []group.MemberIndex{3, 5, 7},
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: nil,
					Group:          initalizeGroup(groupSize, nil, []group.MemberIndex{3, 5, 7}),
				}
			},
		},
		"more than half of threshold disqualified members - failure": {
			disqualifiedMemberIDs: []group.MemberIndex{3, 5, 7},
			expectedResult: func(groupPublicKey *bn256.G2) *Result {
				return &Result{
					GroupPublicKey: nil,
					Group:          initalizeGroup(groupSize, []group.MemberIndex{3, 5, 7}, []group.MemberIndex{}),
				}
			},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			members, err := initializeFinalizingMembersGroup(threshold, groupSize)
			if err != nil {
				t.Fatal(err)
			}

			groupPubKey := members[0].groupPublicKey
			expectedResult := test.expectedResult(groupPubKey)

			for _, member := range members {
				for _, dq := range test.disqualifiedMemberIDs {
					member.group.MarkMemberAsDisqualified(dq)
				}
				for _, ia := range test.inactiveMemberIDs {
					member.group.MarkMemberAsInactive(ia)
				}

				resultToPublish := member.Result()

				if expectedResult.GroupPublicKey != resultToPublish.GroupPublicKey {
					t.Fatalf(
						"Unexpected group public key\nExpected: [%v]\nActual: [%v]\n",
						expectedResult.GroupPublicKey,
						resultToPublish.GroupPublicKey,
					)
				}
				if !reflect.DeepEqual(expectedResult.Group, resultToPublish.Group) {
					t.Fatalf(
						"Unexpected group information\nExpected: [%v]\nActual:   [%v]\n",
						expectedResult.Group,
						resultToPublish.Group,
					)
				}
			}
		})
	}
}

func initializeFinalizingMembersGroup(threshold, groupSize int) ([]*FinalizingMember, error) {
	combiningMembers, err := initializeCombiningMembersGroup(threshold, groupSize)
	if err != nil {
		return nil, err
	}

	var finalizingMembers []*FinalizingMember
	for _, cm := range combiningMembers {
		finalizingMembers = append(finalizingMembers, cm.InitializeFinalization())
	}
	return finalizingMembers, nil
}
