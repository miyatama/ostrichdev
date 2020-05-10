package main

import (
	"testing"
)

func TestHasArgsError(t *testing.T) {
	t.Run("repository is empty", func(t *testing.T) {
		repository := ""
		fromBrancch := ""
		commitID := ""
		ostrichBranch := ""
		err := HasArgsError(&repository, &fromBrancch, &commitID, &ostrichBranch)
		if err == nil {
			t.Fatal("can not get error.")
		}
	})
	t.Run("from branch is empty", func(t *testing.T) {
		repository := "http://miyata.com"
		fromBrancch := ""
		commitID := ""
		ostrichBranch := ""
		err := HasArgsError(&repository, &fromBrancch, &commitID, &ostrichBranch)
		if err == nil {
			t.Fatal("can not get error.")
		}
	})
	t.Run("commit id is empty", func(t *testing.T) {
		repository := "http://miyata.com"
		fromBrancch := "master"
		commitID := ""
		ostrichBranch := ""
		err := HasArgsError(&repository, &fromBrancch, &commitID, &ostrichBranch)
		if err == nil {
			t.Fatal("can not get error.")
		}
	})
	t.Run("ostrich branch is empty", func(t *testing.T) {
		repository := "http://miyata.com"
		fromBrancch := "master"
		commitID := "kfj;alkefja"
		ostrichBranch := ""
		err := HasArgsError(&repository, &fromBrancch, &commitID, &ostrichBranch)
		if err == nil {
			t.Fatal("can not get error.")
		}
	})
	t.Run("all cleear", func(t *testing.T) {
		repository := "http://miyata.com"
		fromBrancch := "master"
		commitID := "kfj;alkefja"
		ostrichBranch := "ostrich"
		err := HasArgsError(&repository, &fromBrancch, &commitID, &ostrichBranch)
		if err != nil {
			t.Fatalf("return error.%#v", err)
		}
	})
}
