package jd

import (
	"sort"
)

type jsonSet jsonArray

var _ JsonNode = jsonSet(nil)

func (s jsonSet) Json() string {
	sMap := make(map[[8]byte]JsonNode)
	for _, n := range s {
		hc := n.hashCode()
		sMap[hc] = n
	}
	hashes := make(hashCodes, 0, len(sMap))
	for hc := range sMap {
		hashes = append(hashes, hc)
	}
	sort.Sort(hashes)
	set := make(jsonSet, 0, len(sMap))
	for _, hc := range hashes {
		set = append(set, sMap[hc])
	}
	return renderJson(set)
}

func (s1 jsonSet) Equals(n JsonNode) bool {
	s2, ok := n.(jsonSet)
	if !ok {
		return false
	}
	if s1.hashCode() == s2.hashCode() {
		return true
	} else {
		return false
	}
}

func (s jsonSet) hashCode() [8]byte {
	sMap := make(map[[8]byte]bool)
	for _, v := range s {
		hc := v.hashCode()
		sMap[hc] = true
	}
	hashes := make(hashCodes, 0, len(sMap))
	for hc := range sMap {
		hashes = append(hashes, hc)
	}
	sort.Sort(hashes)
	b := make([]byte, 0, len(hashes)*8)
	for _, hc := range hashes {
		b = append(b, hc[:]...)
	}
	return hash(b)
}

func (s jsonSet) Diff(n JsonNode) Diff {
	return s.diff(n, Path{})
}

func (s1 jsonSet) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	s2, ok := n.(jsonSet)
	if !ok {
		// Different types
		e := DiffElement{
			Path:      path.clone(),
			OldValues: nodeList(s1),
			NewValues: nodeList(n),
		}
		return append(d, e)
	}
	s1Map := make(map[[8]byte]JsonNode)
	for _, v := range s1 {
		hc := v.hashCode()
		s1Map[hc] = v
	}
	s2Map := make(map[[8]byte]JsonNode)
	for _, v := range s2 {
		hc := v.hashCode()
		s2Map[hc] = v
	}
	s1Hashes := make(hashCodes, 0)
	for hc := range s1Map {
		s1Hashes = append(s1Hashes, hc)
	}
	sort.Sort(s1Hashes)
	s2Hashes := make(hashCodes, 0)
	for hc := range s2Map {
		s2Hashes = append(s2Hashes, hc)
	}
	sort.Sort(s2Hashes)
	e := DiffElement{
		Path:      append(path.clone(), map[string]interface{}{}),
		OldValues: nodeList(),
		NewValues: nodeList(),
	}
	for _, hc := range s1Hashes {
		_, ok := s2Map[hc]
		if !ok {
			e.OldValues = append(e.OldValues, s1Map[hc])
		}
	}
	for _, hc := range s2Hashes {
		_, ok := s1Map[hc]
		if !ok {
			e.NewValues = append(e.NewValues, s2Map[hc])
		}
	}
	if len(e.OldValues) > 0 || len(e.NewValues) > 0 {
		d = append(d, e)
	}
	return d
}

func (s jsonSet) Patch(d Diff) (JsonNode, error) {
	return patchAll(s, d)
}

func (s jsonSet) patch(pathBehind, pathAhead Path, oldValues, newValues []JsonNode) (JsonNode, error) {
	return nil, nil
}
