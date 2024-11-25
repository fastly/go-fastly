package fastly

import (
	"net/http"

	fsterr "github.com/fastly/go-fastly/v9/pkg/errors"
)

// The contents of this file are intended to provide backwards
// compatibility for when the 'errors' package above was part of the
// 'fastly' package itself

type HTTPError = fsterr.HTTPError

func NewHTTPError(resp *http.Response) *HTTPError {
	return fsterr.NewHTTPError(resp)
}

var ErrCommonNameNotInDomains = fsterr.ErrCommonNameNotInDomains
var ErrInvalidMethod = fsterr.ErrInvalidMethod
var ErrManagedLoggingEnabled = fsterr.ErrManagedLoggingEnabled
var ErrMaxExceededEntries = fsterr.ErrMaxExceededEntries
var ErrMaxExceededItems = fsterr.ErrMaxExceededItems
var ErrMaxExceededRules = fsterr.ErrMaxExceededRules
var ErrMissingACLID = fsterr.ErrMissingACLID
var ErrMissingBackend = fsterr.ErrMissingBackend
var ErrMissingCertBlob = fsterr.ErrMissingCertBlob
var ErrMissingCertBundle = fsterr.ErrMissingCertBundle
var ErrMissingCertificateMTLS = fsterr.ErrMissingCertificateMTLS
var ErrMissingCustomerID = fsterr.ErrMissingCustomerID
var ErrMissingDictionaryID = fsterr.ErrMissingDictionaryID
var ErrMissingDirector = fsterr.ErrMissingDirector
var ErrMissingERLID = fsterr.ErrMissingERLID
var ErrMissingEntryID = fsterr.ErrMissingEntryID
var ErrMissingEventID = fsterr.ErrMissingEventID
var ErrMissingFrom = fsterr.ErrMissingFrom
var ErrMissingID = fsterr.ErrMissingID
var ErrMissingImageOptimizerDefaultSetting = fsterr.ErrMissingImageOptimizerDefaultSetting
var ErrMissingIntegrationID = fsterr.ErrMissingIntegrationID
var ErrMissingIntermediatesBlob = fsterr.ErrMissingIntermediatesBlob
var ErrMissingItemKey = fsterr.ErrMissingItemKey
var ErrMissingKey = fsterr.ErrMissingKey
var ErrMissingKeys = fsterr.ErrMissingKeys
var ErrMissingKind = fsterr.ErrMissingKind
var ErrMissingLogin = fsterr.ErrMissingLogin
var ErrMissingMonth = fsterr.ErrMissingMonth
var ErrMissingName = fsterr.ErrMissingName
var ErrMissingNumber = fsterr.ErrMissingNumber
var ErrMissingPermission = fsterr.ErrMissingPermission
var ErrMissingPoolID = fsterr.ErrMissingPoolID
var ErrMissingProductID = fsterr.ErrMissingProductID
var ErrMissingResourceID = fsterr.ErrMissingResourceID
var ErrMissingSecret = fsterr.ErrMissingSecret
var ErrMissingServer = fsterr.ErrMissingServer
var ErrMissingServerSideEncryptionKMSKeyID = fsterr.ErrMissingServerSideEncryptionKMSKeyID
var ErrMissingServiceAuthorizationsService = fsterr.ErrMissingServiceAuthorizationsService
var ErrMissingServiceAuthorizationsUser = fsterr.ErrMissingServiceAuthorizationsUser
var ErrMissingServiceID = fsterr.ErrMissingServiceID
var ErrMissingServiceVersion = fsterr.ErrMissingServiceVersion
var ErrMissingSnippetID = fsterr.ErrMissingSnippetID
var ErrMissingStoreID = fsterr.ErrMissingStoreID
var ErrMissingTLSCertificate = fsterr.ErrMissingTLSCertificate
var ErrMissingTLSDomain = fsterr.ErrMissingTLSDomain
var ErrMissingTo = fsterr.ErrMissingTo
var ErrMissingTokenID = fsterr.ErrMissingTokenID
var ErrMissingTokensValue = fsterr.ErrMissingTokensValue
var ErrMissingURL = fsterr.ErrMissingURL
var ErrMissingUserID = fsterr.ErrMissingUserID
var ErrMissingWAFActiveRule = fsterr.ErrMissingWAFActiveRule
var ErrMissingWAFID = fsterr.ErrMissingWAFID
var ErrMissingWAFRuleExclusion = fsterr.ErrMissingWAFRuleExclusion
var ErrMissingWAFRuleExclusionNumber = fsterr.ErrMissingWAFRuleExclusionNumber
var ErrMissingWAFVersionID = fsterr.ErrMissingWAFVersionID
var ErrMissingWAFVersionNumber = fsterr.ErrMissingWAFVersionNumber
var ErrMissingYear = fsterr.ErrMissingYear
var ErrNotImplemented = fsterr.ErrNotImplemented
var ErrNotOK = fsterr.ErrNotOK
var ErrStatusNotOk = fsterr.ErrStatusNotOk
