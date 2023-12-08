// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//

package consumer

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/http"
	"sort"

	"github.com/antihax/optional"
	nrf_cache "github.com/omec-project/nrf/nrfcache"
	"github.com/omec-project/openapi/Nnrf_NFDiscovery"
	"github.com/omec-project/openapi/models"
	pcf_context "github.com/omec-project/pcf/context"
	"github.com/omec-project/pcf/logger"
	"github.com/omec-project/pcf/util"
)

var SendSearchNFInstances = func(
	nrfUri string, targetNfType, requestNfType models.NfType, param Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (
	models.SearchResult, error) {
	if pcf_context.PCF_Self().EnableNrfCaching {
		logger.ConsumerLog.Errorf("Retrieving from NRF cache")
		return nrf_cache.SearchNFInstances(nrfUri, targetNfType, requestNfType, &param)
	} else {
		return SendNfDiscoveryToNrf(nrfUri, targetNfType, requestNfType, &param)
	}
}

/*
var SendSearchNFInstances = func(
	nrfUri string, targetNfType, requestNfType models.NfType, param Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (
	*models.SearchResult, error) {
	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	result, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requestNfType, &param)
	if err != nil {
		logger.ConsumerLog.Errorf("SearchNFInstances failed: %+v", err)
	}
	defer func() {
		if res != nil {
			if resCloseErr := res.Body.Close(); resCloseErr != nil {
				logger.ConsumerLog.Errorf("NFInstancesStoreApi response body cannot close: %+v", resCloseErr)
			}
		}
	}()
	if res != nil && res.StatusCode == http.StatusTemporaryRedirect {
		return nil, fmt.Errorf("Temporary Redirect For Non NRF Consumer")
	}

	return &result, nil
}*/

func SendNFIntancesUDR(nrfUri, id string) string {
	targetNfType := models.NfType_UDR
	requestNfType := models.NfType_PCF
	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		// 	DataSet: optional.NewInterface(models.DataSetId_SUBSCRIPTION),
	}
	// switch types {
	// case NFDiscoveryToUDRParamSupi:
	// 	localVarOptionals.Supi = optional.NewString(id)
	// case NFDiscoveryToUDRParamExtGroupId:
	// 	localVarOptionals.ExternalGroupIdentity = optional.NewString(id)
	// case NFDiscoveryToUDRParamGpsi:
	// 	localVarOptionals.Gpsi = optional.NewString(id)
	// }

	result, err := SendSearchNFInstances(nrfUri, targetNfType, requestNfType, localVarOptionals)
	if err != nil {
		logger.Consumerlog.Error(err.Error())
		return ""
	}
	nfInstanceIds := make([]string, 0, len(result.NfInstances))
	for _, profile := range result.NfInstances {
		nfInstanceIds = append(nfInstanceIds, profile.NfInstanceId)
	}
	sort.Strings(nfInstanceIds)

	nfInstanceIdIndexMap := make(map[string]int)
	for index, value := range nfInstanceIds {
		nfInstanceIdIndexMap[value] = index
	}

	h := fnv.New32a()
	h.Write([]byte(id))
	key := int(h.Sum32())
	nfInstanceIndex := 0
	if pcf_context.PCF_Self().EnableScaling && len(result.NfInstances) > 0 {
		nfInstanceIndex = key % len(result.NfInstances)
	}

	for _, profile := range result.NfInstances {
		if nfInstanceIndex != nfInstanceIdIndexMap[profile.NfInstanceId] {
			continue
		}
		if uri := util.SearchNFServiceUri(profile, models.ServiceName_NUDR_DR, models.NfServiceStatus_REGISTERED); uri != "" {
			fmt.Println("for Ue: ", id, " for targetNfType ", string(targetNfType), " NF is: ", profile)
			return uri
		}
	}
	return ""
}

func SendNFIntancesAMF(nrfUri string, guami models.Guami, serviceName models.ServiceName, ipAddress string) string {
	targetNfType := models.NfType_AMF
	requestNfType := models.NfType_PCF

	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		Guami: optional.NewInterface(util.MarshToJsonString(guami)),
	}
	// switch types {
	// case NFDiscoveryToUDRParamSupi:
	// 	localVarOptionals.Supi = optional.NewString(id)
	// case NFDiscoveryToUDRParamExtGroupId:
	// 	localVarOptionals.ExternalGroupIdentity = optional.NewString(id)
	// case NFDiscoveryToUDRParamGpsi:
	// 	localVarOptionals.Gpsi = optional.NewString(id)
	// }

	result, err := SendSearchNFInstances(nrfUri, targetNfType, requestNfType, localVarOptionals)
	if err != nil {
		logger.Consumerlog.Error(err.Error())
		return ""
	}

	for _, profile := range result.NfInstances {
		//logger.ConsumerLog.Infof("Looking for ip match %v %v %d", profile.Ipv4Addresses[0], ipAddress, len(result.NfInstances))
		if profile.Ipv4Addresses[0] != ipAddress {
			continue
		}
		return util.SearchNFServiceUri(profile, serviceName, models.NfServiceStatus_REGISTERED)
	}
	return ""
}

func SendNfDiscoveryToNrf(nrfUri string, targetNfType, requestNfType models.NfType,
	param *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (models.SearchResult, error) {
	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	logger.ConsumerLog.Infof("Inside SendNFDiscoveryToNrf")
	result, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requestNfType, param)
	if res != nil && res.StatusCode == http.StatusTemporaryRedirect {
		err = fmt.Errorf("Temporary Redirect For Non NRF Consumer")
	}
	defer func() {
		if bodyCloseErr := res.Body.Close(); bodyCloseErr != nil {
			err = fmt.Errorf("SearchNFInstances' response body cannot close: %+w", bodyCloseErr)
		}
	}()

	pcfSelf := pcf_context.PCF_Self()

	for _, nfProfile := range result.NfInstances {
		// checking whether the PCF subscribed to this target nfinstanceid or not
		logger.ConsumerLog.Infof("Rcvd data from NRF %s %s", nfProfile.NfInstanceId, nfProfile.NfType)
		if _, ok := pcfSelf.NfStatusSubscriptions.Load(nfProfile.NfInstanceId); !ok {
			nrfSubscriptionData := models.NrfSubscriptionData{
				NfStatusNotificationUri: fmt.Sprintf("%s/npcf-callback/v1/nf-status-notify", pcfSelf.GetIPv4Uri()),
				SubscrCond:              &models.NfInstanceIdCond{NfInstanceId: nfProfile.NfInstanceId},
				ReqNfType:               requestNfType,
			}
			nrfSubData, problemDetails, err := SendCreateSubscription(nrfUri, nrfSubscriptionData)
			if problemDetails != nil {
				logger.ConsumerLog.Errorf("SendCreateSubscription to NRF, Problem[%+v]", problemDetails)
			} else if err != nil {
				logger.ConsumerLog.Errorf("SendCreateSubscription Error[%+v]", err)
			}
			logger.ConsumerLog.Infof("Storing subscription %s %s %s", pcfSelf.GetIPv4Uri(), nfProfile.NfInstanceId, nrfSubData.SubscriptionId)
			pcfSelf.NfStatusSubscriptions.Store(nfProfile.NfInstanceId, nrfSubData.SubscriptionId)
		}
	}

	return result, err
}
