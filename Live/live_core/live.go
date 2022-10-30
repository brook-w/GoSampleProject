package live_core

import (
	"Live/config"
	"crypto/md5"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"time"
)

type LiveInterface interface {
	GetPushUrl(streamName string) map[string]string
	GetPullUrl(streamName string) map[string]string
	GetLiveStreamOnlineList(pageNum, pageSize uint64, streamName ...string) (*live.DescribeLiveStreamOnlineListResponse, error)
	ForbidLiveStream(streamName string, resumeTime time.Time, reason string) (*live.ForbidLiveStreamResponse, error)
	DropLiveStream(streamName string)
	ResumeLiveStream(streamName string)
	GetLiveStreamState(streamName string)
	GetForbidLiveStream(pageNum, pageSize uint64, streamName ...string)
}

var lv LiveInterface

func NewLive() LiveInterface {
	ts := config.Conf.GetInt("LiveConf.Ts")
	if lv != nil {
		return lv
	}
	lv = &Live{
		AppName:    config.Conf.GetString("LiveConf.AppName"),
		PushDomain: config.Conf.GetString("LiveConf.PushDomain"),
		PushKey:    config.Secret.GetString("Live.PushKey"),
		PullDomain: config.Conf.GetString("LiveConf.PullDomain"),
		PullKey:    config.Secret.GetString("Live.PullKey"),
		Ts:         time.Second * time.Duration(ts),
	}
	return lv
}

// LiveInterface
// 顺序接口
// 1. 获取推流地址（鉴权）
// 2. 获取播放地址（鉴权）
// 3. 查询直播中的流
// 4. 禁推流
// 5. 断开流
// 6. 恢复直播流
// 7. 查询流状态
// 8. 获取禁推流列表

type Live struct {
	AppName    string
	PushDomain string
	PushKey    string
	PullDomain string
	PullKey    string
	Ts         time.Duration
}

func (l *Live) DropLiveStream(streamName string) {
	//TODO implement me
	panic("implement me")
}

func (l *Live) ResumeLiveStream(streamName string) {
	//TODO implement me
	panic("implement me")
}

func (l *Live) GetLiveStreamState(streamName string) {
	//TODO implement me
	panic("implement me")
}

func (l *Live) GetForbidLiveStream(pageNum, pageSize uint64, streamName ...string) {
	//TODO implement me
	panic("implement me")
}

func (l *Live) GetPushUrl(streamName string) map[string]string {
	mp := map[string]string{
		"rtmp":   fmt.Sprintf("rtmp://%s/%s/%s", l.PushDomain, l.AppName, streamName),
		"webrtc": fmt.Sprintf("webrtc://%s/%s/%s", l.PushDomain, l.AppName, streamName),
		"srt":    fmt.Sprintf("srt://%s:9000?streamId=#!::h=/%s,r=%s/%s", l.PushDomain, l.PushDomain, l.AppName, streamName),
	}

	if l.PushKey != "" {
		txSecret, txTime := sign(l.PushKey, streamName, l.Ts)
		for k, v := range mp {
			if k == "srt" {
				mp[k] = fmt.Sprintf("%s,txSecret=%s,txTime=%s", v, txSecret, txTime)
			} else {
				mp[k] = fmt.Sprintf("%s?txSecret=%s&txTime=%s", v, txSecret, txTime)
			}
		}
	}
	return mp
}

func (l *Live) GetPullUrl(streamName string) map[string]string {
	mp := map[string]string{
		"rtmp":   fmt.Sprintf("rtmp://%s/%s/%s", l.PullDomain, l.AppName, streamName),
		"webrtc": fmt.Sprintf("webrtc://%s/%s/%s", l.PullDomain, l.AppName, streamName),
		"srt":    fmt.Sprintf("srt://%s:9000?streamId=#!::h=/%s,r=%s/%s", l.PullDomain, l.PushDomain, l.AppName, streamName),
		"hls":    fmt.Sprintf("http://%s/%s/%s.m3u8", l.PullDomain, l.AppName, streamName),
		"flv":    fmt.Sprintf("http://%s/%s/%s.flv", l.PullDomain, l.AppName, streamName),
	}

	if l.PullKey != "" {
		txSecret, txTime := sign(l.PullKey, streamName, l.Ts)
		for k, v := range mp {
			mp[k] = fmt.Sprintf("%s?txSecret=%s&txTime=%s", v, txSecret, txTime)
		}
	}
	return mp
}

func (l *Live) GetLiveStreamOnlineList(pageNum, pageSize uint64, streamName ...string) (*live.DescribeLiveStreamOnlineListResponse, error) {
	credential := getCredential()
	client := getClient(credential)

	request := live.NewDescribeLiveStreamOnlineListRequest()
	request.DomainName = common.StringPtr(l.PushDomain)
	request.AppName = common.StringPtr(l.AppName)
	request.PageNum = common.Uint64Ptr(pageNum)
	request.PageSize = common.Uint64Ptr(pageSize)
	if len(streamName) > 0 {
		request.StreamName = common.StringPtr(streamName[0])
	}

	// 返回的resp是一个DescribeLiveStreamOnlineListResponse的实例，与请求对象对应
	response, err := client.DescribeLiveStreamOnlineList(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	return response, nil
}

func (l *Live) ForbidLiveStream(streamName string, resumeTime time.Time, reason string) (*live.ForbidLiveStreamResponse, error) {
	credential := getCredential()
	client := getClient(credential)

	request := live.NewForbidLiveStreamRequest()
	request.AppName = common.StringPtr(l.AppName)
	request.DomainName = common.StringPtr(l.PushDomain)
	request.StreamName = common.StringPtr(streamName)
	if resumeTime.IsZero() {
		request.ResumeTime = common.StringPtr(resumeTime.UTC().String())
	}
	request.Reason = common.StringPtr(reason)

	response, err := client.ForbidLiveStream(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}

	return response, nil

}

func getClient(credential *common.Credential) *live.Client {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = config.Conf.GetString("LiveConf.EndPoint")
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := live.NewClient(credential, "", cpf)
	return client
}

func getCredential() *common.Credential {
	return common.NewCredential(
		config.Secret.GetString("TencentCloudSecret.SecretId"),
		config.Secret.GetString("TencentCloudSecret.SecretKey"),
	)
}

func sign(key, streamName string, duration time.Duration) (txSecret, txTime string) {
	endTime := time.Now().Add(duration).Unix()
	txTime = fmt.Sprintf("%x", endTime)
	bytes := md5.Sum([]byte(fmt.Sprintf("%s%s%s", key, streamName, txTime)))
	txSecret = fmt.Sprintf("%x", bytes)
	return
}
