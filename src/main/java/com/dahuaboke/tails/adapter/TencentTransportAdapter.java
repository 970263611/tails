package com.dahuaboke.tails.adapter;

import com.dahuaboke.tails.util.CommonUtils;

import java.io.UnsupportedEncodingException;
import java.net.URLEncoder;
import java.util.TreeMap;

/**
 * author: dahua
 * date: 2023/12/13 16:20
 */
public class TencentTransportAdapter extends TransportAdapter {

    private static final String URL = "https://fanyi-api.baidu.com/api/trans/vip/translate?";

    private String appid;
    private String securityKey;

    public TencentTransportAdapter(String appid, String securityKey) {
        this.appid = appid;
        this.securityKey = securityKey;
    }

    @Override
    protected String doTransport(String text) {
        try {
            TreeMap<String, Object> params = new TreeMap<String, Object>(); // TreeMap可以自动排序
            // 实际调用时应当使用随机数，例如：params.put("Nonce", new Random().nextInt(java.lang.Integer.MAX_VALUE));
            params.put("Nonce", 11886); // 公共参数
            // 实际调用时应当使用系统当前时间，例如：   params.put("Timestamp", System.currentTimeMillis() / 1000);
            params.put("Timestamp", 1465185768); // 公共参数
            // 需要设置环境变量 TENCENTCLOUD_SECRET_ID，值为示例的 AKIDz8krbsJ5yKBZQpn74WFkmLPx3*******
            params.put("SecretId", System.getenv("TENCENTCLOUD_SECRET_ID")); // 公共参数
            params.put("Action", "DescribeInstances"); // 公共参数
            params.put("Version", "2017-03-12"); // 公共参数
            params.put("Region", "ap-guangzhou"); // 公共参数
            params.put("Limit", 20); // 业务参数
            params.put("Offset", 0); // 业务参数
            params.put("InstanceIds.0", "ins-09dx96dg"); // 业务参数
            // 需要设置环境变量 TENCENTCLOUD_SECRET_KEY，值为示例的 Gu5t9xGARNpq86cd98joQYCN3*******
            params.put("Signature", CommonUtils.sign(getStringToSign(params), System.getenv("TENCENTCLOUD_SECRET_KEY"), "HmacSHA1")); // 公共参数
            return getUrl(params);
        } catch (Exception e) {
            return null;
        }
    }

    public static String getStringToSign(TreeMap<String, Object> params) {
        StringBuilder s2s = new StringBuilder("GETcvm.tencentcloudapi.com/?");
        // 签名时要求对参数进行字典排序，此处用TreeMap保证顺序
        for (String k : params.keySet()) {
            s2s.append(k).append("=").append(params.get(k).toString()).append("&");
        }
        return s2s.toString().substring(0, s2s.length() - 1);
    }

    public static String getUrl(TreeMap<String, Object> params) throws UnsupportedEncodingException {
        StringBuilder url = new StringBuilder("https://cvm.tencentcloudapi.com/?");
        // 实际请求的url中对参数顺序没有要求
        for (String k : params.keySet()) {
            // 需要对请求串进行urlencode，由于key都是英文字母，故此处仅对其value进行urlencode
            url.append(k).append("=").append(URLEncoder.encode(params.get(k).toString(), "utf-8")).append("&");
        }
        return url.toString().substring(0, url.length() - 1);
    }
}
