package yqlib

import (
	"container/list"

	"gopkg.in/yaml.v3"
)

func createMapOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	log.Debugf("-- createMapOperation")

	//each matchingNodes entry should turn into a sequence of keys to create.
	//then collect object should do a cross function of the same index sequence for all matches.

	var path []interface{}

	var document uint = 0

	sequences := list.New()

	if matchingNodes.Len() > 0 {

		for matchingNodeEl := matchingNodes.Front(); matchingNodeEl != nil; matchingNodeEl = matchingNodeEl.Next() {
			matchingNode := matchingNodeEl.Value.(*CandidateNode)
			sequenceNode, err := sequenceFor(d, matchingNode, pathNode)
			if err != nil {
				return nil, err
			}
			sequences.PushBack(sequenceNode)
		}
	} else {
		sequenceNode, err := sequenceFor(d, nil, pathNode)
		if err != nil {
			return nil, err
		}
		sequences.PushBack(sequenceNode)
	}

	return nodeToMap(&CandidateNode{Node: listToNodeSeq(sequences), Document: document, Path: path}), nil

}

func sequenceFor(d *dataTreeNavigator, matchingNode *CandidateNode, pathNode *PathTreeNode) (*CandidateNode, error) {
	var path []interface{}
	var document uint = 0
	var matches = list.New()

	if matchingNode != nil {
		path = matchingNode.Path
		document = matchingNode.Document
		matches = nodeToMap(matchingNode)
	}

	mapPairs, err := crossFunction(d, matches, pathNode,
		func(d *dataTreeNavigator, lhs *CandidateNode, rhs *CandidateNode) (*CandidateNode, error) {
			node := yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			log.Debugf("LHS:", NodeToString(lhs))
			log.Debugf("RHS:", NodeToString(rhs))
			node.Content = []*yaml.Node{
				UnwrapDoc(lhs.Node),
				UnwrapDoc(rhs.Node),
			}

			return &CandidateNode{Node: &node, Document: document, Path: path}, nil
		})

	if err != nil {
		return nil, err
	}
	innerList := listToNodeSeq(mapPairs)
	innerList.Style = yaml.FlowStyle
	return &CandidateNode{Node: innerList, Document: document, Path: path}, nil
}

//NOTE: here the document index gets dropped so we
// no longer know where the node originates from.
func listToNodeSeq(list *list.List) *yaml.Node {
	node := yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for entry := list.Front(); entry != nil; entry = entry.Next() {
		entryCandidate := entry.Value.(*CandidateNode)
		log.Debugf("Collecting %v into sequence", NodeToString(entryCandidate))
		node.Content = append(node.Content, entryCandidate.Node)
	}
	return &node
}
