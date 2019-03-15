// Copyright 2018-2019 opcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gopcua/opcua"
	uid "github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"
)

type NodeDef struct {
	NodeID      *ua.NodeID
	Path        string
	DataType    string
	Writable    bool
	Description string
}

func join(a, b string) string {
	if a == "" {
		return b
	}
	return a + "." + b
}

func browse(n *opcua.Node, path string, level int) ([]NodeDef, error) {
	if level > 10 {
		return nil, nil
	}
	nodeClass, err := n.NodeClass()
	if err != nil {
		return nil, err
	}
	browseName, err := n.BrowseName()
	if err != nil {
		return nil, err
	}
	descr, err := n.Description()
	if err != nil {
		return nil, err
	}
	accessLevel, err := n.AccessLevel()
	if err != nil {
		return nil, err
	}
	path = join(path, browseName.Name)

	switch nodeClass {
	case ua.NodeClassObject:
		var nodes []NodeDef
		children, err := n.Children()
		if err != nil {
			return nil, err
		}
		for _, cn := range children {
			childnodes, err := browse(cn, path, level+1)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, childnodes...)
		}
		return nodes, nil

	case ua.NodeClassVariable:
		return []NodeDef{
			{
				NodeID:      n.ID,
				Path:        path,
				Description: descr.Text,
				Writable:    accessLevel&0x2 == 0x2,
				DataType:    "???",
			},
		}, nil

	}

	typeDefs := ua.NewTwoByteNodeID(uid.HasTypeDefinition)
	refs, err := n.References(typeDefs)
	if err != nil {
		return nil, err
	}
	// todo(fs): example still incomplete
	// log.Printf("refs: %#v err: %v", refs, err)
	for _, r := range refs.Results {
		for _, ref := range r.References {
			log.Printf("resp: %s, %#v", ref.NodeID, ref.BrowseName)
		}
	}
	return nil, nil
}

func main() {
	var (
		endpoint = flag.String("endpoint", "opc.tcp://localhost:4840", "OPC UA Endpoint URL")
		nodeID   = flag.String("node", "", "node id for the root node")
	)
	flag.Parse()

	id, err := ua.NewNodeID(*nodeID)
	if err != nil {
		log.Fatalf("invalid node id: %s", err)
	}

	c := opcua.NewClient(*endpoint, nil)
	if err := c.Open(); err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	nodeList, err := browse(c.Node(id), "", 0)
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range nodeList {
		fmt.Println(s)
	}
}
