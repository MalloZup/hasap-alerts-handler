package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"testing"
)

func TestHanaPrimaryTrue(t *testing.T) {
	xmlFile, err := os.Open("test/crm_mon.xml")
	if err != nil {
		t.Error(err)
	}

	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var cMon CrmMon

	xml.Unmarshal(byteValue, &cMon)
	hanaRes := lookUpHanaNodePrimary(&cMon, "hana01")
	if hanaRes != "msl_SAPHana_PRD_HDB00" {
		t.Logf("GOT: %s, EXPECTED: msl_SAPHana_PRD_HDB00", hanaRes)
		t.Errorf("HanaNodePrimary should have returned the resource id")
	}
}

func TestHanaPrimaryFalse(t *testing.T) {
	xmlFile, err := os.Open("test/crm_mon.xml")
	if err != nil {
		t.Error(err)
	}

	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var cMon CrmMon

	xml.Unmarshal(byteValue, &cMon)
	if lookUpHanaNodePrimary(&cMon, "hana02") != "" {
		t.Errorf("we should not have found primary node")
	}
}
