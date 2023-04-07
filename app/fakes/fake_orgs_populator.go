package fakes

import (
	"fmt"
	"os"
	"path"
	"ticken-event-service/infra"
	"ticken-event-service/models"
)

type SeedOrganization struct {
	AdminUsername string `json:"admin_username"`
	Name          string `json:"name"`
}

func (loader *Loader) seedOrganizations(toSeed []*SeedOrganization) []error {
	var seedErrors = make([]error, 0)

	if loader.repoProvider.GetOrganizationRepository().Count() > 0 {
		return seedErrors
	}

	for _, organization := range toSeed {
		// load genesis org in the database

		organizer := loader.repoProvider.GetOrganizerRepository().FindOrganizerByUsername(organization.AdminUsername)
		if organizer == nil {
			seedErrors = append(
				seedErrors,
				fmt.Errorf("failed to seed organization %s: organizer with username %s not found", organization.Name, organization.AdminUsername),
			)
		}

		newOrg := createOrganization(
			organization.Name,
			organizer,
			loader.config.Pvtbc.ClusterStoragePath,
			loader.hsm,
		)

		if err := loader.repoProvider.GetOrganizationRepository().AddOne(newOrg); err != nil {
			panic(err)
		}
	}
	return seedErrors
}

func createOrganization(orgName string, admin *models.Organizer, clusterStoragePath string, hsm infra.HSM) *models.Organization {
	orgCACert, tlsCACert := readOrgMSP(orgName, clusterStoragePath)
	newOrganization := models.NewOrganization(orgName, getOrgChannelName(orgName), orgCACert, tlsCACert)

	userOrgCert := readOrgUserMSP(orgName, admin.Username, clusterStoragePath, hsm)
	_ = newOrganization.AddUser(admin, "admin", userOrgCert)

	nodeOrgCert, nodeTlsCert := readOrgNodeMSP(orgName, "peer0", clusterStoragePath, hsm)
	_ = newOrganization.AddNode("peer0", orgName+"-peer0"+".localho.st:443", nodeOrgCert, nodeTlsCert)

	return newOrganization
}

func readOrgMSP(orgName string, clusterStoragePath string) (*models.Certificate, *models.Certificate) {
	baseMspPath := path.Join(clusterStoragePath, "orgs", "peer-orgs", orgName, "msp")

	orgCACertBytes, err := os.ReadFile(path.Join(baseMspPath, "cacerts", "ca-signcert.pem"))
	if err != nil {
		panic(err)
	}

	tlsCACertBytes, err := os.ReadFile(path.Join(baseMspPath, "tlscacerts", "tlsca-signcert.pem"))
	if err != nil {
		panic(err)
	}

	return models.NewCertificate(orgCACertBytes, ""), models.NewCertificate(tlsCACertBytes, "")
}

func readOrgUserMSP(orgName string, _ string, clusterStoragePath string, hsm infra.HSM) *models.Certificate {
	// todo -> replace here and in the pvtbc bootstrap the admin name
	userMspPath := path.Join(clusterStoragePath, "orgs", "peer-orgs", orgName, "users", orgName+"-admin", "msp")

	userCertBytes, err := os.ReadFile(path.Join(userMspPath, "signcerts", "cert.pem"))
	if err != nil {
		panic(err)
	}

	userPrivBytes, err := os.ReadFile(path.Join(userMspPath, "keystore", "priv.pem"))
	if err != nil {
		panic(err)
	}

	userPrivStorageKey, err := hsm.Store(userPrivBytes)
	if err != nil {
		panic(err)
	}

	return models.NewCertificate(userCertBytes, userPrivStorageKey)
}

func readOrgNodeMSP(orgName string, node string, clusterStoragePath string, hsm infra.HSM) (*models.Certificate, *models.Certificate) {
	nodePath := path.Join(clusterStoragePath, "orgs", "peer-orgs", orgName, "nodes", orgName+"-"+node)

	nodeCertBytes, err := os.ReadFile(path.Join(nodePath, "msp", "signcerts", "cert.pem"))
	if err != nil {
		panic(err)
	}

	nodePrivBytes, err := os.ReadFile(path.Join(nodePath, "msp", "keystore", "priv.pem"))
	if err != nil {
		panic(err)
	}

	nodePrivStorageKey, err := hsm.Store(nodePrivBytes)
	if err != nil {
		panic(err)
	}

	tlsCertBytes, err := os.ReadFile(path.Join(nodePath, "tls", "signcerts", "tls-cert.pem"))
	if err != nil {
		panic(err)
	}

	return models.NewCertificate(nodeCertBytes, nodePrivStorageKey), models.NewCertificate(tlsCertBytes, "")
}

func getOrgChannelName(orgName string) string {
	return "ticken" + "-" + orgName + "-" + "channel"
}
