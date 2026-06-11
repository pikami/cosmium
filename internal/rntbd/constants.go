package rntbd

import (
	"fmt"

	"github.com/pikami/cosmium/api/headers"
)

type RntbdOperationType uint16

const (
	RntbdOperationTypeConnection                       RntbdOperationType = 0x0000
	RntbdOperationTypeCreate                           RntbdOperationType = 0x0001
	RntbdOperationTypeUpdate                           RntbdOperationType = 0x0002
	RntbdOperationTypeRead                             RntbdOperationType = 0x0003
	RntbdOperationTypeReadFeed                         RntbdOperationType = 0x0004
	RntbdOperationTypeDelete                           RntbdOperationType = 0x0005
	RntbdOperationTypeReplace                          RntbdOperationType = 0x0006
	RntbdOperationTypeExecuteJavaScript                RntbdOperationType = 0x0008
	RntbdOperationTypeSQLQuery                         RntbdOperationType = 0x0009
	RntbdOperationTypePause                            RntbdOperationType = 0x000A
	RntbdOperationTypeResume                           RntbdOperationType = 0x000B
	RntbdOperationTypeStop                             RntbdOperationType = 0x000C
	RntbdOperationTypeRecycle                          RntbdOperationType = 0x000D
	RntbdOperationTypeCrash                            RntbdOperationType = 0x000E
	RntbdOperationTypeQuery                            RntbdOperationType = 0x000F
	RntbdOperationTypeForceConfigRefresh               RntbdOperationType = 0x0010
	RntbdOperationTypeHead                             RntbdOperationType = 0x0011
	RntbdOperationTypeHeadFeed                         RntbdOperationType = 0x0012
	RntbdOperationTypeUpsert                           RntbdOperationType = 0x0013
	RntbdOperationTypeRecreate                         RntbdOperationType = 0x0014
	RntbdOperationTypeThrottle                         RntbdOperationType = 0x0015
	RntbdOperationTypeGetSplitPoint                    RntbdOperationType = 0x0016
	RntbdOperationTypePreCreateValidation              RntbdOperationType = 0x0017
	RntbdOperationTypeBatchApply                       RntbdOperationType = 0x0018
	RntbdOperationTypeAbortSplit                       RntbdOperationType = 0x0019
	RntbdOperationTypeCompleteSplit                    RntbdOperationType = 0x001A
	RntbdOperationTypeOfferUpdateOperation             RntbdOperationType = 0x001B
	RntbdOperationTypeOfferPreGrowValidation           RntbdOperationType = 0x001C
	RntbdOperationTypeBatchReportThroughputUtilization RntbdOperationType = 0x001D
	RntbdOperationTypeCompletePartitionMigration       RntbdOperationType = 0x001E
	RntbdOperationTypeAbortPartitionMigration          RntbdOperationType = 0x001F
	RntbdOperationTypePreReplaceValidation             RntbdOperationType = 0x0020
	RntbdOperationTypeAddComputeGatewayRequestCharges  RntbdOperationType = 0x0021
	RntbdOperationTypeMigratePartition                 RntbdOperationType = 0x0022
)

type RntbdResourceType uint16

const (
	RntbdResourceTypeConnection              RntbdResourceType = 0x0000
	RntbdResourceTypeDatabase                RntbdResourceType = 0x0001
	RntbdResourceTypeCollection              RntbdResourceType = 0x0002
	RntbdResourceTypeDocument                RntbdResourceType = 0x0003
	RntbdResourceTypeAttachment              RntbdResourceType = 0x0004
	RntbdResourceTypeUser                    RntbdResourceType = 0x0005
	RntbdResourceTypePermission              RntbdResourceType = 0x0006
	RntbdResourceTypeStoredProcedure         RntbdResourceType = 0x0007
	RntbdResourceTypeConflict                RntbdResourceType = 0x0008
	RntbdResourceTypeTrigger                 RntbdResourceType = 0x0009
	RntbdResourceTypeUserDefinedFunction     RntbdResourceType = 0x000A
	RntbdResourceTypeModule                  RntbdResourceType = 0x000B
	RntbdResourceTypeReplica                 RntbdResourceType = 0x000C
	RntbdResourceTypeModuleCommand           RntbdResourceType = 0x000D
	RntbdResourceTypeRecord                  RntbdResourceType = 0x000E
	RntbdResourceTypeOffer                   RntbdResourceType = 0x000F
	RntbdResourceTypePartitionSetInformation RntbdResourceType = 0x0010
	RntbdResourceTypeXPReplicatorAddress     RntbdResourceType = 0x0011
	RntbdResourceTypeMasterPartition         RntbdResourceType = 0x0012
	RntbdResourceTypeServerPartition         RntbdResourceType = 0x0013
	RntbdResourceTypeDatabaseAccount         RntbdResourceType = 0x0014
	RntbdResourceTypeTopology                RntbdResourceType = 0x0015
	RntbdResourceTypePartitionKeyRange       RntbdResourceType = 0x0016
	RntbdResourceTypeSchema                  RntbdResourceType = 0x0018
	RntbdResourceTypeBatchApply              RntbdResourceType = 0x0019
	RntbdResourceTypeRestoreMetadata         RntbdResourceType = 0x001A
	RntbdResourceTypeComputeGatewayCharges   RntbdResourceType = 0x001B
	RntbdResourceTypeRidRange                RntbdResourceType = 0x001C
	RntbdResourceTypeUserDefinedType         RntbdResourceType = 0x001D
)

type RntbdRequestHeader uint16

const (
	RntbdRequestHeaderResourceId                                RntbdRequestHeader = 0x0000 // RntbdTokenType.Bytes, required = false
	RntbdRequestHeaderAuthorizationToken                        RntbdRequestHeader = 0x0001 // RntbdTokenType.String, required = false
	RntbdRequestHeaderPayloadPresent                            RntbdRequestHeader = 0x0002 // RntbdTokenType.Byte, required = true
	RntbdRequestHeaderDate                                      RntbdRequestHeader = 0x0003 // RntbdTokenType.SmallString, required = false
	RntbdRequestHeaderPageSize                                  RntbdRequestHeader = 0x0004 // RntbdTokenType.ULong, required = false
	RntbdRequestHeaderSessionToken                              RntbdRequestHeader = 0x0005 // RntbdTokenType.String, required = false
	RntbdRequestHeaderContinuationToken                         RntbdRequestHeader = 0x0006 // RntbdTokenType.String, required = false
	RntbdRequestHeaderIndexingDirective                         RntbdRequestHeader = 0x0007 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderMatch                                     RntbdRequestHeader = 0x0008 // RntbdTokenType.String, required = false
	RntbdRequestHeaderPreTriggerInclude                         RntbdRequestHeader = 0x0009 // RntbdTokenType.String, required = false
	RntbdRequestHeaderPostTriggerInclude                        RntbdRequestHeader = 0x000A // RntbdTokenType.String, required = false
	RntbdRequestHeaderIsFanout                                  RntbdRequestHeader = 0x000B // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderCollectionPartitionIndex                  RntbdRequestHeader = 0x000C // RntbdTokenType.ULong, required = false
	RntbdRequestHeaderCollectionServiceIndex                    RntbdRequestHeader = 0x000D // RntbdTokenType.ULong, required = false
	RntbdRequestHeaderPreTriggerExclude                         RntbdRequestHeader = 0x000E // RntbdTokenType.String, required = false
	RntbdRequestHeaderPostTriggerExclude                        RntbdRequestHeader = 0x000F // RntbdTokenType.String, required = false
	RntbdRequestHeaderConsistencyLevel                          RntbdRequestHeader = 0x0010 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderEntityId                                  RntbdRequestHeader = 0x0011 // RntbdTokenType.String, required = false
	RntbdRequestHeaderResourceSchemaName                        RntbdRequestHeader = 0x0012 // RntbdTokenType.SmallString, required = false
	RntbdRequestHeaderReplicaPath                               RntbdRequestHeader = 0x0013 // RntbdTokenType.String, required = true
	RntbdRequestHeaderResourceTokenExpiry                       RntbdRequestHeader = 0x0014 // RntbdTokenType.ULong, required = false
	RntbdRequestHeaderDatabaseName                              RntbdRequestHeader = 0x0015 // RntbdTokenType.String, required = false
	RntbdRequestHeaderCollectionName                            RntbdRequestHeader = 0x0016 // RntbdTokenType.String, required = false
	RntbdRequestHeaderDocumentName                              RntbdRequestHeader = 0x0017 // RntbdTokenType.String, required = false
	RntbdRequestHeaderAttachmentName                            RntbdRequestHeader = 0x0018 // RntbdTokenType.String, required = false
	RntbdRequestHeaderUserName                                  RntbdRequestHeader = 0x0019 // RntbdTokenType.String, required = false
	RntbdRequestHeaderPermissionName                            RntbdRequestHeader = 0x001A // RntbdTokenType.String, required = false
	RntbdRequestHeaderStoredProcedureName                       RntbdRequestHeader = 0x001B // RntbdTokenType.String, required = false
	RntbdRequestHeaderUserDefinedFunctionName                   RntbdRequestHeader = 0x001C // RntbdTokenType.String, required = false
	RntbdRequestHeaderTriggerName                               RntbdRequestHeader = 0x001D // RntbdTokenType.String, required = false
	RntbdRequestHeaderEnableScanInQuery                         RntbdRequestHeader = 0x001E // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderEmitVerboseTracesInQuery                  RntbdRequestHeader = 0x001F // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderConflictName                              RntbdRequestHeader = 0x0020 // RntbdTokenType.String, required = false
	RntbdRequestHeaderBindReplicaDirective                      RntbdRequestHeader = 0x0021 // RntbdTokenType.String, required = false
	RntbdRequestHeaderPrimaryMasterKey                          RntbdRequestHeader = 0x0022 // RntbdTokenType.String, required = false
	RntbdRequestHeaderSecondaryMasterKey                        RntbdRequestHeader = 0x0023 // RntbdTokenType.String, required = false
	RntbdRequestHeaderPrimaryReadonlyKey                        RntbdRequestHeader = 0x0024 // RntbdTokenType.String, required = false
	RntbdRequestHeaderSecondaryReadonlyKey                      RntbdRequestHeader = 0x0025 // RntbdTokenType.String, required = false
	RntbdRequestHeaderProfileRequest                            RntbdRequestHeader = 0x0026 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderEnableLowPrecisionOrderBy                 RntbdRequestHeader = 0x0027 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderClientVersion                             RntbdRequestHeader = 0x0028 // RntbdTokenType.SmallString, required = false
	RntbdRequestHeaderCanCharge                                 RntbdRequestHeader = 0x0029 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderCanThrottle                               RntbdRequestHeader = 0x002A // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderPartitionKey                              RntbdRequestHeader = 0x002B // RntbdTokenType.String, required = false
	RntbdRequestHeaderPartitionKeyRangeId                       RntbdRequestHeader = 0x002C // RntbdTokenType.String, required = false
	RntbdRequestHeaderNotUsed2D                                 RntbdRequestHeader = 0x002D // RntbdTokenType.Invalid, required = false
	RntbdRequestHeaderNotUsed2E                                 RntbdRequestHeader = 0x002E // RntbdTokenType.Invalid, required = false
	RntbdRequestHeaderNotUsed2F                                 RntbdRequestHeader = 0x002F // RntbdTokenType.Invalid, required = false
	RntbdRequestHeaderMigrateCollectionDirective                RntbdRequestHeader = 0x0031 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderNotUsed32                                 RntbdRequestHeader = 0x0032 // RntbdTokenType.Invalid, required = false
	RntbdRequestHeaderSupportSpatialLegacyCoordinates           RntbdRequestHeader = 0x0033 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderPartitionCount                            RntbdRequestHeader = 0x0034 // RntbdTokenType.ULong, required = false
	RntbdRequestHeaderCollectionRid                             RntbdRequestHeader = 0x0035 // RntbdTokenType.String, required = false
	RntbdRequestHeaderPartitionKeyRangeName                     RntbdRequestHeader = 0x0036 // RntbdTokenType.String, required = false
	RntbdRequestHeaderSchemaName                                RntbdRequestHeader = 0x003A // RntbdTokenType.String, required = false
	RntbdRequestHeaderFilterBySchemaRid                         RntbdRequestHeader = 0x003B // RntbdTokenType.String, required = false
	RntbdRequestHeaderUsePolygonsSmallerThanAHemisphere         RntbdRequestHeader = 0x003C // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderGatewaySignature                          RntbdRequestHeader = 0x003D // RntbdTokenType.String, required = false
	RntbdRequestHeaderEnableLogging                             RntbdRequestHeader = 0x003E // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderAIM                                       RntbdRequestHeader = 0x003F // RntbdTokenType.String, required = false
	RntbdRequestHeaderPopulateQuotaInfo                         RntbdRequestHeader = 0x0040 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderDisableRUPerMinuteUsage                   RntbdRequestHeader = 0x0041 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderPopulateQueryMetrics                      RntbdRequestHeader = 0x0042 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderResponseContinuationTokenLimitInKb        RntbdRequestHeader = 0x0043 // RntbdTokenType.ULong, required = false
	RntbdRequestHeaderPopulatePartitionStatistics               RntbdRequestHeader = 0x0044 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderRemoteStorageType                         RntbdRequestHeader = 0x0045 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderCollectionRemoteStorageSecurityIdentifier RntbdRequestHeader = 0x0046 // RntbdTokenType.String, required = false
	RntbdRequestHeaderIfModifiedSince                           RntbdRequestHeader = 0x0047 // RntbdTokenType.String, required = false
	RntbdRequestHeaderPopulateCollectionThroughputInfo          RntbdRequestHeader = 0x0048 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderRemainingTimeInMsOnClientRequest          RntbdRequestHeader = 0x0049 // RntbdTokenType.ULong, required = false
	RntbdRequestHeaderClientRetryAttemptCount                   RntbdRequestHeader = 0x004A // RntbdTokenType.ULong, required = false
	RntbdRequestHeaderTargetLsn                                 RntbdRequestHeader = 0x004B // RntbdTokenType.LongLong, required = false
	RntbdRequestHeaderTargetGlobalCommittedLsn                  RntbdRequestHeader = 0x004C // RntbdTokenType.LongLong, required = false
	RntbdRequestHeaderTransportRequestID                        RntbdRequestHeader = 0x004D // RntbdTokenType.ULong, required = false
	RntbdRequestHeaderRestoreMetadaFilter                       RntbdRequestHeader = 0x004E // RntbdTokenType.String, required = false
	RntbdRequestHeaderRestoreParams                             RntbdRequestHeader = 0x004F // RntbdTokenType.String, required = false
	RntbdRequestHeaderShareThroughput                           RntbdRequestHeader = 0x0050 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderPartitionResourceFilter                   RntbdRequestHeader = 0x0051 // RntbdTokenType.String, required = false
	RntbdRequestHeaderIsReadOnlyScript                          RntbdRequestHeader = 0x0052 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderIsAutoScaleRequest                        RntbdRequestHeader = 0x0053 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderForceQueryScan                            RntbdRequestHeader = 0x0054 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderCanOfferReplaceComplete                   RntbdRequestHeader = 0x0056 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderExcludeSystemProperties                   RntbdRequestHeader = 0x0057 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderBinaryId                                  RntbdRequestHeader = 0x0058 // RntbdTokenType.Bytes, required = false
	RntbdRequestHeaderTimeToLiveInSeconds                       RntbdRequestHeader = 0x0059 // RntbdTokenType.Long, required = false
	RntbdRequestHeaderEffectivePartitionKey                     RntbdRequestHeader = 0x005A // RntbdTokenType.Bytes, required = false
	RntbdRequestHeaderBinaryPassthroughRequest                  RntbdRequestHeader = 0x005B // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderUserDefinedTypeName                       RntbdRequestHeader = 0x005C // RntbdTokenType.String, required = false
	RntbdRequestHeaderEnableDynamicRidRangeAllocation           RntbdRequestHeader = 0x005D // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderEnumerationDirection                      RntbdRequestHeader = 0x005E // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderStartId                                   RntbdRequestHeader = 0x005F // RntbdTokenType.Bytes, required = false
	RntbdRequestHeaderEndId                                     RntbdRequestHeader = 0x0060 // RntbdTokenType.Bytes, required = false
	RntbdRequestHeaderFanoutOperationState                      RntbdRequestHeader = 0x0061 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderStartEpk                                  RntbdRequestHeader = 0x0062 // RntbdTokenType.Bytes, required = false
	RntbdRequestHeaderEndEpk                                    RntbdRequestHeader = 0x0063 // RntbdTokenType.Bytes, required = false
	RntbdRequestHeaderReadFeedKeyType                           RntbdRequestHeader = 0x0064 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderContentSerializationFormat                RntbdRequestHeader = 0x0065 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderAllowTentativeWrites                      RntbdRequestHeader = 0x0066 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderIsUserRequest                             RntbdRequestHeader = 0x0067 // RntbdTokenType.Byte, required = false
	RntbdRequestHeaderSharedOfferThroughput                     RntbdRequestHeader = 0x0068 // RntbdTokenType.ULong, required = false

	RntbdRequestHeaderSDKSupportedCapabilities RntbdRequestHeader = 0x00A2 // RntbdTokenType.ULong, required = ?
)

type RntbdResponseHeaderType uint16

const (
	RntbdResponseHeaderPayloadPresent               RntbdResponseHeaderType = 0x0000 // RntbdTokenType.Byte, required = true
	RntbdResponseHeaderLastStateChangeDateTime      RntbdResponseHeaderType = 0x0002 // RntbdTokenType.SmallString, required = false
	RntbdResponseHeaderContinuationToken            RntbdResponseHeaderType = 0x0003 // RntbdTokenType.String, required = false
	RntbdResponseHeaderETag                         RntbdResponseHeaderType = 0x0004 // RntbdTokenType.String, required = false
	RntbdResponseHeaderReadsPerformed               RntbdResponseHeaderType = 0x0007 // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderWritesPerformed              RntbdResponseHeaderType = 0x0008 // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderQueriesPerformed             RntbdResponseHeaderType = 0x0009 // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderIndexTermsGenerated          RntbdResponseHeaderType = 0x000A // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderScriptsExecuted              RntbdResponseHeaderType = 0x000B // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderRetryAfterMilliseconds       RntbdResponseHeaderType = 0x000C // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderIndexingDirective            RntbdResponseHeaderType = 0x000D // RntbdTokenType.Byte, required = false
	RntbdResponseHeaderStorageMaxResoureQuota       RntbdResponseHeaderType = 0x000E // RntbdTokenType.String, required = false
	RntbdResponseHeaderStorageResourceQuotaUsage    RntbdResponseHeaderType = 0x000F // RntbdTokenType.String, required = false
	RntbdResponseHeaderSchemaVersion                RntbdResponseHeaderType = 0x0010 // RntbdTokenType.SmallString, required = false
	RntbdResponseHeaderCollectionPartitionIndex     RntbdResponseHeaderType = 0x0011 // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderCollectionServiceIndex       RntbdResponseHeaderType = 0x0012 // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderLSN                          RntbdResponseHeaderType = 0x0013 // RntbdTokenType.LongLong, required = false
	RntbdResponseHeaderItemCount                    RntbdResponseHeaderType = 0x0014 // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderRequestCharge                RntbdResponseHeaderType = 0x0015 // RntbdTokenType.Double, required = false
	RntbdResponseHeaderOwnerFullName                RntbdResponseHeaderType = 0x0017 // RntbdTokenType.String, required = false
	RntbdResponseHeaderOwnerId                      RntbdResponseHeaderType = 0x0018 // RntbdTokenType.String, required = false
	RntbdResponseHeaderDatabaseAccountId            RntbdResponseHeaderType = 0x0019 // RntbdTokenType.String, required = false
	RntbdResponseHeaderQuorumAckedLSN               RntbdResponseHeaderType = 0x001A // RntbdTokenType.LongLong, required = false
	RntbdResponseHeaderRequestValidationFailure     RntbdResponseHeaderType = 0x001B // RntbdTokenType.Byte, required = false
	RntbdResponseHeaderSubStatus                    RntbdResponseHeaderType = 0x001C // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderCollectionUpdateProgress     RntbdResponseHeaderType = 0x001D // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderCurrentWriteQuorum           RntbdResponseHeaderType = 0x001E // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderCurrentReplicaSetSize        RntbdResponseHeaderType = 0x001F // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderCollectionLazyIndexProgress  RntbdResponseHeaderType = 0x0020 // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderPartitionKeyRangeId          RntbdResponseHeaderType = 0x0021 // RntbdTokenType.String, required = false
	RntbdResponseHeaderLogResults                   RntbdResponseHeaderType = 0x0025 // RntbdTokenType.String, required = false
	RntbdResponseHeaderXPRole                       RntbdResponseHeaderType = 0x0026 // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderIsRUPerMinuteUsed            RntbdResponseHeaderType = 0x0027 // RntbdTokenType.Byte, required = false
	RntbdResponseHeaderQueryMetrics                 RntbdResponseHeaderType = 0x0028 // RntbdTokenType.String, required = false
	RntbdResponseHeaderGlobalCommittedLSN           RntbdResponseHeaderType = 0x0029 // RntbdTokenType.LongLong, required = false
	RntbdResponseHeaderNumberOfReadRegions          RntbdResponseHeaderType = 0x0030 // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderOfferReplacePending          RntbdResponseHeaderType = 0x0031 // RntbdTokenType.Byte, required = false
	RntbdResponseHeaderItemLSN                      RntbdResponseHeaderType = 0x0032 // RntbdTokenType.LongLong, required = false
	RntbdResponseHeaderRestoreState                 RntbdResponseHeaderType = 0x0033 // RntbdTokenType.String, required = false
	RntbdResponseHeaderCollectionSecurityIdentifier RntbdResponseHeaderType = 0x0034 // RntbdTokenType.String, required = false
	RntbdResponseHeaderTransportRequestID           RntbdResponseHeaderType = 0x0035 // RntbdTokenType.ULong, required = false
	RntbdResponseHeaderShareThroughput              RntbdResponseHeaderType = 0x0036 // RntbdTokenType.Byte, required = false
	RntbdResponseHeaderDisableRntbdChannel          RntbdResponseHeaderType = 0x0038 // RntbdTokenType.Byte, required = false
	RntbdResponseHeaderServerDateTimeUtc            RntbdResponseHeaderType = 0x0039 // RntbdTokenType.SmallString, required = false
	RntbdResponseHeaderLocalLSN                     RntbdResponseHeaderType = 0x003A // RntbdTokenType.LongLong, required = false
	RntbdResponseHeaderQuorumAckedLocalLSN          RntbdResponseHeaderType = 0x003B // RntbdTokenType.LongLong, required = false
	RntbdResponseHeaderItemLocalLSN                 RntbdResponseHeaderType = 0x003C // RntbdTokenType.LongLong, required = false
	RntbdResponseHeaderHasTentativeWrites           RntbdResponseHeaderType = 0x003D // RntbdTokenType.Byte, required = false
	RntbdResponseHeaderSessionToken                 RntbdResponseHeaderType = 0x003E // RntbdTokenType.String, required = false
)

type RntbdContextHeader uint16

const (
	RntbdContextHeaderProtocolVersion                 RntbdContextHeader = 0x0000 // RntbdTokenType.ULong, required = false
	RntbdContextHeaderClientVersion                   RntbdContextHeader = 0x0001 // RntbdTokenType.SmallString, required = false
	RntbdContextHeaderServerAgent                     RntbdContextHeader = 0x0002 // RntbdTokenType.SmallString, required = true
	RntbdContextHeaderServerVersion                   RntbdContextHeader = 0x0003 // RntbdTokenType.SmallString, required = true
	RntbdContextHeaderIdleTimeoutInSeconds            RntbdContextHeader = 0x0004 // RntbdTokenType.ULong, required = false
	RntbdContextHeaderUnauthenticatedTimeoutInSeconds RntbdContextHeader = 0x0005 // RntbdTokenType.ULong, required = false
)

type RntbdTokenType uint8

const (
	RntbdTokenTypeByte        RntbdTokenType = 0x00 // 8bit boolean
	RntbdTokenTypeUShort      RntbdTokenType = 0x01 // 16bit unsigned integer
	RntbdTokenTypeULong       RntbdTokenType = 0x02 // 32bit unsigned integer
	RntbdTokenTypeLong        RntbdTokenType = 0x03 // 32bit signed integer
	RntbdTokenTypeULongLong   RntbdTokenType = 0x04 // 64bit unsigned integer
	RntbdTokenTypeLongLong    RntbdTokenType = 0x05 // 64bit signed integer
	RntbdTokenTypeGuid        RntbdTokenType = 0x06 // 128bit GUID
	RntbdTokenTypeSmallString RntbdTokenType = 0x07 // 8bit len + string
	RntbdTokenTypeString      RntbdTokenType = 0x08 // 16bit len + string
	RntbdTokenTypeULongString RntbdTokenType = 0x09 // 32bit len + string
	RntbdTokenTypeSmallBytes  RntbdTokenType = 0x0A // 8bit len + bytes
	RntbdTokenTypeBytes       RntbdTokenType = 0x0B // 16bit len + bytes
	RntbdTokenTypeULongBytes  RntbdTokenType = 0x0C // 32bit len + bytes
	RntbdTokenTypeFloat       RntbdTokenType = 0x0D // 32bit float
	RntbdTokenTypeDouble      RntbdTokenType = 0x0E // 64bit double
	RntbdTokenTypeInvalid     RntbdTokenType = 0x0F // Invalid token type
)

func (h RntbdRequestHeader) String() string {
	switch h {
	case RntbdRequestHeaderResourceId:
		return "RntbdRequestHeaderResourceId"
	case RntbdRequestHeaderAuthorizationToken:
		return headers.Authorization
	case RntbdRequestHeaderPayloadPresent:
		return "RntbdRequestHeaderPayloadPresent"
	case RntbdRequestHeaderDate:
		return headers.XDate
	case RntbdRequestHeaderPageSize:
		return "RntbdRequestHeaderPageSize"
	case RntbdRequestHeaderSessionToken:
		return "RntbdRequestHeaderSessionToken"
	case RntbdRequestHeaderContinuationToken:
		return "RntbdRequestHeaderContinuationToken"
	case RntbdRequestHeaderIndexingDirective:
		return "RntbdRequestHeaderIndexingDirective"
	case RntbdRequestHeaderMatch:
		return "RntbdRequestHeaderMatch"
	case RntbdRequestHeaderPreTriggerInclude:
		return "RntbdRequestHeaderPreTriggerInclude"
	case RntbdRequestHeaderPostTriggerInclude:
		return "RntbdRequestHeaderPostTriggerInclude"
	case RntbdRequestHeaderIsFanout:
		return "RntbdRequestHeaderIsFanout"
	case RntbdRequestHeaderCollectionPartitionIndex:
		return "RntbdRequestHeaderCollectionPartitionIndex"
	case RntbdRequestHeaderCollectionServiceIndex:
		return "RntbdRequestHeaderCollectionServiceIndex"
	case RntbdRequestHeaderPreTriggerExclude:
		return "RntbdRequestHeaderPreTriggerExclude"
	case RntbdRequestHeaderPostTriggerExclude:
		return "RntbdRequestHeaderPostTriggerExclude"
	case RntbdRequestHeaderConsistencyLevel:
		return headers.ConsistencyLevel
	case RntbdRequestHeaderEntityId:
		return "RntbdRequestHeaderEntityId"
	case RntbdRequestHeaderResourceSchemaName:
		return "RntbdRequestHeaderResourceSchemaName"
	case RntbdRequestHeaderReplicaPath:
		return "RntbdRequestHeaderReplicaPath"
	case RntbdRequestHeaderResourceTokenExpiry:
		return "RntbdRequestHeaderResourceTokenExpiry"
	case RntbdRequestHeaderDatabaseName:
		return "RntbdRequestHeaderDatabaseName"
	case RntbdRequestHeaderCollectionName:
		return "RntbdRequestHeaderCollectionName"
	case RntbdRequestHeaderDocumentName:
		return "RntbdRequestHeaderDocumentName"
	case RntbdRequestHeaderAttachmentName:
		return "RntbdRequestHeaderAttachmentName"
	case RntbdRequestHeaderUserName:
		return "RntbdRequestHeaderUserName"
	case RntbdRequestHeaderPermissionName:
		return "RntbdRequestHeaderPermissionName"
	case RntbdRequestHeaderStoredProcedureName:
		return "RntbdRequestHeaderStoredProcedureName"
	case RntbdRequestHeaderUserDefinedFunctionName:
		return "RntbdRequestHeaderUserDefinedFunctionName"
	case RntbdRequestHeaderTriggerName:
		return "RntbdRequestHeaderTriggerName"
	case RntbdRequestHeaderEnableScanInQuery:
		return "RntbdRequestHeaderEnableScanInQuery"
	case RntbdRequestHeaderEmitVerboseTracesInQuery:
		return "RntbdRequestHeaderEmitVerboseTracesInQuery"
	case RntbdRequestHeaderConflictName:
		return "RntbdRequestHeaderConflictName"
	case RntbdRequestHeaderBindReplicaDirective:
		return "RntbdRequestHeaderBindReplicaDirective"
	case RntbdRequestHeaderPrimaryMasterKey:
		return "RntbdRequestHeaderPrimaryMasterKey"
	case RntbdRequestHeaderSecondaryMasterKey:
		return "RntbdRequestHeaderSecondaryMasterKey"
	case RntbdRequestHeaderPrimaryReadonlyKey:
		return "RntbdRequestHeaderPrimaryReadonlyKey"
	case RntbdRequestHeaderSecondaryReadonlyKey:
		return "RntbdRequestHeaderSecondaryReadonlyKey"
	case RntbdRequestHeaderProfileRequest:
		return "RntbdRequestHeaderProfileRequest"
	case RntbdRequestHeaderEnableLowPrecisionOrderBy:
		return "RntbdRequestHeaderEnableLowPrecisionOrderBy"
	case RntbdRequestHeaderClientVersion:
		return "RntbdRequestHeaderClientVersion"
	case RntbdRequestHeaderCanCharge:
		return "RntbdRequestHeaderCanCharge"
	case RntbdRequestHeaderCanThrottle:
		return "RntbdRequestHeaderCanThrottle"
	case RntbdRequestHeaderPartitionKey:
		return "RntbdRequestHeaderPartitionKey"
	case RntbdRequestHeaderPartitionKeyRangeId:
		return "RntbdRequestHeaderPartitionKeyRangeId"
	case RntbdRequestHeaderNotUsed2D:
		return "RntbdRequestHeaderNotUsed2D"
	case RntbdRequestHeaderNotUsed2E:
		return "RntbdRequestHeaderNotUsed2E"
	case RntbdRequestHeaderNotUsed2F:
		return "RntbdRequestHeaderNotUsed2F"
	case RntbdRequestHeaderMigrateCollectionDirective:
		return "RntbdRequestHeaderMigrateCollectionDirective"
	case RntbdRequestHeaderNotUsed32:
		return "RntbdRequestHeaderNotUsed32"
	case RntbdRequestHeaderSupportSpatialLegacyCoordinates:
		return "RntbdRequestHeaderSupportSpatialLegacyCoordinates"
	case RntbdRequestHeaderPartitionCount:
		return "RntbdRequestHeaderPartitionCount"
	case RntbdRequestHeaderCollectionRid:
		return "RntbdRequestHeaderCollectionRid"
	case RntbdRequestHeaderPartitionKeyRangeName:
		return "RntbdRequestHeaderPartitionKeyRangeName"
	case RntbdRequestHeaderSchemaName:
		return "RntbdRequestHeaderSchemaName"
	case RntbdRequestHeaderFilterBySchemaRid:
		return "RntbdRequestHeaderFilterBySchemaRid"
	case RntbdRequestHeaderUsePolygonsSmallerThanAHemisphere:
		return "RntbdRequestHeaderUsePolygonsSmallerThanAHemisphere"
	case RntbdRequestHeaderGatewaySignature:
		return "RntbdRequestHeaderGatewaySignature"
	case RntbdRequestHeaderEnableLogging:
		return "RntbdRequestHeaderEnableLogging"
	case RntbdRequestHeaderAIM:
		return headers.AIM
	case RntbdRequestHeaderPopulateQuotaInfo:
		return "RntbdRequestHeaderPopulateQuotaInfo"
	case RntbdRequestHeaderDisableRUPerMinuteUsage:
		return "RntbdRequestHeaderDisableRUPerMinuteUsage"
	case RntbdRequestHeaderPopulateQueryMetrics:
		return "RntbdRequestHeaderPopulateQueryMetrics"
	case RntbdRequestHeaderResponseContinuationTokenLimitInKb:
		return "RntbdRequestHeaderResponseContinuationTokenLimitInKb"
	case RntbdRequestHeaderPopulatePartitionStatistics:
		return "RntbdRequestHeaderPopulatePartitionStatistics"
	case RntbdRequestHeaderRemoteStorageType:
		return "RntbdRequestHeaderRemoteStorageType"
	case RntbdRequestHeaderCollectionRemoteStorageSecurityIdentifier:
		return "RntbdRequestHeaderCollectionRemoteStorageSecurityIdentifier"
	case RntbdRequestHeaderIfModifiedSince:
		return "RntbdRequestHeaderIfModifiedSince"
	case RntbdRequestHeaderPopulateCollectionThroughputInfo:
		return "RntbdRequestHeaderPopulateCollectionThroughputInfo"
	case RntbdRequestHeaderRemainingTimeInMsOnClientRequest:
		return headers.RemainingTimeInMsOnClient
	case RntbdRequestHeaderClientRetryAttemptCount:
		return headers.ClientRetryAttemptCount
	case RntbdRequestHeaderTargetLsn:
		return "RntbdRequestHeaderTargetLsn"
	case RntbdRequestHeaderTargetGlobalCommittedLsn:
		return "RntbdRequestHeaderTargetGlobalCommittedLsn"
	case RntbdRequestHeaderTransportRequestID:
		return "RntbdRequestHeaderTransportRequestID"
	case RntbdRequestHeaderRestoreMetadaFilter:
		return "RntbdRequestHeaderRestoreMetadaFilter"
	case RntbdRequestHeaderRestoreParams:
		return "RntbdRequestHeaderRestoreParams"
	case RntbdRequestHeaderShareThroughput:
		return "RntbdRequestHeaderShareThroughput"
	case RntbdRequestHeaderPartitionResourceFilter:
		return "RntbdRequestHeaderPartitionResourceFilter"
	case RntbdRequestHeaderIsReadOnlyScript:
		return "RntbdRequestHeaderIsReadOnlyScript"
	case RntbdRequestHeaderIsAutoScaleRequest:
		return "RntbdRequestHeaderIsAutoScaleRequest"
	case RntbdRequestHeaderForceQueryScan:
		return "RntbdRequestHeaderForceQueryScan"
	case RntbdRequestHeaderCanOfferReplaceComplete:
		return "RntbdRequestHeaderCanOfferReplaceComplete"
	case RntbdRequestHeaderExcludeSystemProperties:
		return "RntbdRequestHeaderExcludeSystemProperties"
	case RntbdRequestHeaderBinaryId:
		return "RntbdRequestHeaderBinaryId"
	case RntbdRequestHeaderTimeToLiveInSeconds:
		return "RntbdRequestHeaderTimeToLiveInSeconds"
	case RntbdRequestHeaderEffectivePartitionKey:
		return "RntbdRequestHeaderEffectivePartitionKey"
	case RntbdRequestHeaderBinaryPassthroughRequest:
		return "RntbdRequestHeaderBinaryPassthroughRequest"
	case RntbdRequestHeaderUserDefinedTypeName:
		return "RntbdRequestHeaderUserDefinedTypeName"
	case RntbdRequestHeaderEnableDynamicRidRangeAllocation:
		return "RntbdRequestHeaderEnableDynamicRidRangeAllocation"
	case RntbdRequestHeaderEnumerationDirection:
		return "RntbdRequestHeaderEnumerationDirection"
	case RntbdRequestHeaderStartId:
		return "RntbdRequestHeaderStartId"
	case RntbdRequestHeaderEndId:
		return "RntbdRequestHeaderEndId"
	case RntbdRequestHeaderFanoutOperationState:
		return "RntbdRequestHeaderFanoutOperationState"
	case RntbdRequestHeaderStartEpk:
		return "RntbdRequestHeaderStartEpk"
	case RntbdRequestHeaderEndEpk:
		return "RntbdRequestHeaderEndEpk"
	case RntbdRequestHeaderReadFeedKeyType:
		return "RntbdRequestHeaderReadFeedKeyType"
	case RntbdRequestHeaderContentSerializationFormat:
		return "RntbdRequestHeaderContentSerializationFormat"
	case RntbdRequestHeaderAllowTentativeWrites:
		return "RntbdRequestHeaderAllowTentativeWrites"
	case RntbdRequestHeaderIsUserRequest:
		return "RntbdRequestHeaderIsUserRequest"
	case RntbdRequestHeaderSharedOfferThroughput:
		return "RntbdRequestHeaderSharedOfferThroughput"
	case RntbdRequestHeaderSDKSupportedCapabilities:
		return headers.SupportedCapabilities
	}

	return fmt.Sprintf("RntbdRequestHeader(%d)", h)
}

func (h RntbdContextHeader) String() string {
	switch h {
	case RntbdContextHeaderProtocolVersion:
		return "RntbdContextHeaderProtocolVersion"
	case RntbdContextHeaderClientVersion:
		return "RntbdContextHeaderClientVersion"
	case RntbdContextHeaderServerAgent:
		return "RntbdContextHeaderServerAgent"
	case RntbdContextHeaderServerVersion:
		return "RntbdContextHeaderServerVersion"
	case RntbdContextHeaderIdleTimeoutInSeconds:
		return "RntbdContextHeaderIdleTimeoutInSeconds"
	case RntbdContextHeaderUnauthenticatedTimeoutInSeconds:
		return "RntbdContextHeaderUnauthenticatedTimeoutInSeconds"
	}

	return fmt.Sprintf("RntbdContextHeader(%d)", h)
}

func (h RntbdResponseHeaderType) String() string {
	switch h {
	case RntbdResponseHeaderPayloadPresent:
		return "PayloadPresent"
	case RntbdResponseHeaderLastStateChangeDateTime:
		return "LastStateChangeDateTime"
	case RntbdResponseHeaderContinuationToken:
		return "ContinuationToken"
	case RntbdResponseHeaderETag:
		return "ETag"
	case RntbdResponseHeaderReadsPerformed:
		return "ReadsPerformed"
	case RntbdResponseHeaderWritesPerformed:
		return "WritesPerformed"
	case RntbdResponseHeaderQueriesPerformed:
		return "QueriesPerformed"
	case RntbdResponseHeaderIndexTermsGenerated:
		return "IndexTermsGenerated"
	case RntbdResponseHeaderScriptsExecuted:
		return "ScriptsExecuted"
	case RntbdResponseHeaderRetryAfterMilliseconds:
		return "RetryAfterMilliseconds"
	case RntbdResponseHeaderIndexingDirective:
		return "IndexingDirective"
	case RntbdResponseHeaderStorageMaxResoureQuota:
		return "StorageMaxResoureQuota"
	case RntbdResponseHeaderStorageResourceQuotaUsage:
		return "StorageResourceQuotaUsage"
	case RntbdResponseHeaderSchemaVersion:
		return "SchemaVersion"
	case RntbdResponseHeaderCollectionPartitionIndex:
		return "CollectionPartitionIndex"
	case RntbdResponseHeaderCollectionServiceIndex:
		return "CollectionServiceIndex"
	case RntbdResponseHeaderLSN:
		return "LSN"
	case RntbdResponseHeaderItemCount:
		return "ItemCount"
	case RntbdResponseHeaderRequestCharge:
		return "RequestCharge"
	case RntbdResponseHeaderOwnerFullName:
		return "OwnerFullName"
	case RntbdResponseHeaderOwnerId:
		return "OwnerId"
	case RntbdResponseHeaderDatabaseAccountId:
		return "DatabaseAccountId"
	case RntbdResponseHeaderQuorumAckedLSN:
		return "QuorumAckedLSN"
	case RntbdResponseHeaderRequestValidationFailure:
		return "RequestValidationFailure"
	case RntbdResponseHeaderSubStatus:
		return "SubStatus"
	case RntbdResponseHeaderCollectionUpdateProgress:
		return "CollectionUpdateProgress"
	case RntbdResponseHeaderCurrentWriteQuorum:
		return "CurrentWriteQuorum"
	case RntbdResponseHeaderCurrentReplicaSetSize:
		return "CurrentReplicaSetSize"
	case RntbdResponseHeaderCollectionLazyIndexProgress:
		return "CollectionLazyIndexProgress"
	case RntbdResponseHeaderPartitionKeyRangeId:
		return "PartitionKeyRangeId"
	case RntbdResponseHeaderLogResults:
		return "LogResults"
	case RntbdResponseHeaderXPRole:
		return "XPRole"
	case RntbdResponseHeaderIsRUPerMinuteUsed:
		return "IsRUPerMinuteUsed"
	case RntbdResponseHeaderQueryMetrics:
		return "QueryMetrics"
	case RntbdResponseHeaderGlobalCommittedLSN:
		return "GlobalCommittedLSN"
	case RntbdResponseHeaderNumberOfReadRegions:
		return "NumberOfReadRegions"
	case RntbdResponseHeaderOfferReplacePending:
		return "OfferReplacePending"
	case RntbdResponseHeaderItemLSN:
		return "ItemLSN"
	case RntbdResponseHeaderRestoreState:
		return "RestoreState"
	case RntbdResponseHeaderCollectionSecurityIdentifier:
		return "CollectionSecurityIdentifier"
	case RntbdResponseHeaderTransportRequestID:
		return "TransportRequestID"
	case RntbdResponseHeaderShareThroughput:
		return "ShareThroughput"
	case RntbdResponseHeaderDisableRntbdChannel:
		return "DisableRntbdChannel"
	case RntbdResponseHeaderServerDateTimeUtc:
		return "ServerDateTimeUtc"
	case RntbdResponseHeaderLocalLSN:
		return "LocalLSN"
	case RntbdResponseHeaderQuorumAckedLocalLSN:
		return "QuorumAckedLocalLSN"
	case RntbdResponseHeaderItemLocalLSN:
		return "ItemLocalLSN"
	case RntbdResponseHeaderHasTentativeWrites:
		return "HasTentativeWrites"
	case RntbdResponseHeaderSessionToken:
		return "SessionToken"
	}

	return fmt.Sprintf("RntbdResponseHeaderType(%d)", h)
}

func (r RntbdResourceType) String() string {
	switch r {
	case RntbdResourceTypeConnection:
		return "Connection"
	case RntbdResourceTypeDatabase:
		return "Database"
	case RntbdResourceTypeCollection:
		return "Collection"
	case RntbdResourceTypeDocument:
		return "Document"
	case RntbdResourceTypeAttachment:
		return "Attachment"
	case RntbdResourceTypeUser:
		return "User"
	case RntbdResourceTypePermission:
		return "Permission"
	case RntbdResourceTypeStoredProcedure:
		return "StoredProcedure"
	case RntbdResourceTypeConflict:
		return "Conflict"
	case RntbdResourceTypeTrigger:
		return "Trigger"
	case RntbdResourceTypeUserDefinedFunction:
		return "UserDefinedFunction"
	case RntbdResourceTypeModule:
		return "Module"
	case RntbdResourceTypeReplica:
		return "Replica"
	case RntbdResourceTypeModuleCommand:
		return "ModuleCommand"
	case RntbdResourceTypeRecord:
		return "Record"
	case RntbdResourceTypeOffer:
		return "Offer"
	case RntbdResourceTypePartitionSetInformation:
		return "PartitionSetInformation"
	case RntbdResourceTypeXPReplicatorAddress:
		return "XPReplicatorAddress"
	case RntbdResourceTypeMasterPartition:
		return "MasterPartition"
	case RntbdResourceTypeServerPartition:
		return "ServerPartition"
	case RntbdResourceTypeDatabaseAccount:
		return "DatabaseAccount"
	case RntbdResourceTypeTopology:
		return "Topology"
	case RntbdResourceTypePartitionKeyRange:
		return "PartitionKeyRange"
	case RntbdResourceTypeSchema:
		return "Schema"
	case RntbdResourceTypeBatchApply:
		return "BatchApply"
	case RntbdResourceTypeRestoreMetadata:
		return "RestoreMetadata"
	case RntbdResourceTypeComputeGatewayCharges:
		return "ComputeGatewayCharges"
	case RntbdResourceTypeRidRange:
		return "RidRange"
	case RntbdResourceTypeUserDefinedType:
		return "UserDefinedType"
	}

	return fmt.Sprintf("RntbdResourceType(%d)", r)
}

func (o RntbdOperationType) String() string {
	switch o {
	case RntbdOperationTypeConnection:
		return "Connection"
	case RntbdOperationTypeCreate:
		return "Create"
	case RntbdOperationTypeUpdate:
		return "Update"
	case RntbdOperationTypeRead:
		return "Read"
	case RntbdOperationTypeReadFeed:
		return "ReadFeed"
	case RntbdOperationTypeDelete:
		return "Delete"
	case RntbdOperationTypeReplace:
		return "Replace"
	case RntbdOperationTypeExecuteJavaScript:
		return "ExecuteJavaScript"
	case RntbdOperationTypeSQLQuery:
		return "SQLQuery"
	case RntbdOperationTypePause:
		return "Pause"
	case RntbdOperationTypeResume:
		return "Resume"
	case RntbdOperationTypeStop:
		return "Stop"
	case RntbdOperationTypeRecycle:
		return "Recycle"
	case RntbdOperationTypeCrash:
		return "Crash"
	case RntbdOperationTypeQuery:
		return "Query"
	case RntbdOperationTypeForceConfigRefresh:
		return "ForceConfigRefresh"
	case RntbdOperationTypeHead:
		return "Head"
	case RntbdOperationTypeHeadFeed:
		return "HeadFeed"
	case RntbdOperationTypeUpsert:
		return "Upsert"
	case RntbdOperationTypeRecreate:
		return "Recreate"
	case RntbdOperationTypeThrottle:
		return "Throttle"
	case RntbdOperationTypeGetSplitPoint:
		return "GetSplitPoint"
	case RntbdOperationTypePreCreateValidation:
		return "PreCreateValidation"
	case RntbdOperationTypeBatchApply:
		return "BatchApply"
	case RntbdOperationTypeAbortSplit:
		return "AbortSplit"
	case RntbdOperationTypeCompleteSplit:
		return "CompleteSplit"
	case RntbdOperationTypeOfferUpdateOperation:
		return "OfferUpdateOperation"
	case RntbdOperationTypeOfferPreGrowValidation:
		return "OfferPreGrowValidation"
	case RntbdOperationTypeBatchReportThroughputUtilization:
		return "BatchReportThroughputUtilization"
	case RntbdOperationTypeCompletePartitionMigration:
		return "CompletePartitionMigration"
	case RntbdOperationTypeAbortPartitionMigration:
		return "AbortPartitionMigration"
	case RntbdOperationTypePreReplaceValidation:
		return "PreReplaceValidation"
	case RntbdOperationTypeAddComputeGatewayRequestCharges:
		return "AddComputeGatewayRequestCharges"
	case RntbdOperationTypeMigratePartition:
		return "MigratePartition"
	}

	return fmt.Sprintf("RntbdOperationType(%d)", o)
}
