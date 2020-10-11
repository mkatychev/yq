package treeops

import (
	"github.com/elliotchance/orderedmap"
	"gopkg.in/yaml.v3"
)

type OperatorHandler func(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error)

func TraverseOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	return d.getMatchingNodes(lhs, pathNode.Rhs)
}

func AssignOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	for el := lhs.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		log.Debugf("Assiging %v to %v", node.getKey(), pathNode.Rhs.PathElement.StringValue)
		node.Node.Value = pathNode.Rhs.PathElement.StringValue
	}
	return lhs, nil
}

func UnionOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	rhs, err := d.getMatchingNodes(matchingNodes, pathNode.Rhs)
	if err != nil {
		return nil, err
	}
	for el := rhs.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		lhs.Set(node.getKey(), node)
	}
	return lhs, nil
}

func IntersectionOperator(d *dataTreeNavigator, matchingNodes *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	rhs, err := d.getMatchingNodes(matchingNodes, pathNode.Rhs)
	if err != nil {
		return nil, err
	}
	var matchingNodeMap = orderedmap.NewOrderedMap()
	for el := lhs.Front(); el != nil; el = el.Next() {
		_, exists := rhs.Get(el.Key)
		if exists {
			matchingNodeMap.Set(el.Key, el.Value)
		}
	}
	return matchingNodeMap, nil
}

func splatNode(d *dataTreeNavigator, candidate *CandidateNode) (*orderedmap.OrderedMap, error) {
	elMap := orderedmap.NewOrderedMap()
	elMap.Set(candidate.getKey(), candidate)
	//need to splat matching nodes, then search through them
	splatter := &PathTreeNode{PathElement: &PathElement{
		PathElementType: PathKey,
		Value:           "*",
		StringValue:     "*",
	}}
	return d.getMatchingNodes(elMap, splatter)
}

func EqualsOperator(d *dataTreeNavigator, matchMap *orderedmap.OrderedMap, pathNode *PathTreeNode) (*orderedmap.OrderedMap, error) {
	log.Debugf("-- equalsOperation")
	var results = orderedmap.NewOrderedMap()

	for el := matchMap.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		valuePattern := pathNode.Rhs.PathElement.StringValue
		log.Debug("checking %v", candidate)

		// if pathNode.Lhs.PathElement.PathElementType == SelfReference {
		// 	if Match(candidate.Node.Value, valuePattern) {
		// 		results.Set(el.Key, el.Value)
		// 	}
		// } else {
		errInChild := findMatchingChildren(d, results, candidate, pathNode.Lhs, valuePattern)
		if errInChild != nil {
			return nil, errInChild
		}
		// }

	}

	return results, nil
}

func findMatchingChildren(d *dataTreeNavigator, results *orderedmap.OrderedMap, candidate *CandidateNode, lhs *PathTreeNode, valuePattern string) error {
	var children *orderedmap.OrderedMap
	var err error
	// don't splat scalars.
	if candidate.Node.Kind != yaml.ScalarNode {
		children, err = splatNode(d, candidate)
		log.Debugf("-- splatted matches, ")
		if err != nil {
			return err
		}
	} else {
		children = orderedmap.NewOrderedMap()
		children.Set(candidate.getKey(), candidate)
	}

	for childEl := children.Front(); childEl != nil; childEl = childEl.Next() {
		childMap := orderedmap.NewOrderedMap()
		childMap.Set(childEl.Key, childEl.Value)
		childMatches, errChild := d.getMatchingNodes(childMap, lhs)
		log.Debug("got the LHS")
		if errChild != nil {
			return errChild
		}

		if containsMatchingValue(childMatches, valuePattern) {
			results.Set(childEl.Key, childEl.Value)
		}
	}
	return nil
}

func containsMatchingValue(matchMap *orderedmap.OrderedMap, valuePattern string) bool {
	log.Debugf("-- findMatchingValues")

	for el := matchMap.Front(); el != nil; el = el.Next() {
		node := el.Value.(*CandidateNode)
		log.Debugf("-- compating %v to %v", node.Node.Value, valuePattern)
		if Match(node.Node.Value, valuePattern) {
			return true
		}
	}
	log.Debugf("-- done findMatchingValues")

	return false
}