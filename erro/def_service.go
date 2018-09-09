package erro

import "github.com/czsilence/go/erro"

var (
	E_OAUTH_FAILED              = erro.New(300001, "oauth failed")
	E_OAuth_NoRequestSecret     = erro.New(300002, "no request secret found")
	E_OAuth_OAuthInfoNotFound   = erro.New(300003, "oauth info not found")
	E_OAuth_OAuthConfigNotFound = erro.New(300004, "can not find oauth config")

	E_Transfer_InitTransportFailed = erro.New(300101, "init transport failed")
	E_Transfer_InitProtocolFailed  = erro.New(300102, "init transport protocol failed")
	E_Transfer_InitClientFailed    = erro.New(300103, "init transport client failed")

	E_Transfer_Tag_ListTagsFailed  = erro.New(300201, "list tags failed failed")
	E_Transfer_Tag_CreateTagFailed = erro.New(300202, "create tag failed failed")

	E_Transfer_NB_ListNotebookFailed   = erro.New(300301, "list notebooks failed failed")
	E_Transfer_NB_CreateNotebookFailed = erro.New(300302, "create notebook failed failed")

	E_Transfer_Note_GetSyncChunkFailed = erro.New(300401, "get sync chunk failed failed")
)
