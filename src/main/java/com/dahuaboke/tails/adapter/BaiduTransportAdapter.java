package com.dahuaboke.tails.adapter;

import com.dahuaboke.tails.util.CommonUtils;
import com.dahuaboke.tails.util.HTTP;
import org.apache.commons.lang3.StringUtils;

import java.net.URLEncoder;
import java.util.List;
import java.util.Map;

/**
 * author: dahua
 * date: 2023/12/13 16:20
 */
public class BaiduTransportAdapter extends TransportAdapter {

    private static final String URL = "https://fanyi-api.baidu.com/api/trans/vip/translate?";

    private String appid;
    private String securityKey;

    public BaiduTransportAdapter(String appid, String securityKey) {
        this.appid = appid;
        this.securityKey = securityKey;
    }

    @Override
    protected String doTransport(String text) {
        try {
            String[] split = StringUtils.split(text, "\n");
            StringBuffer stringBuffer = new StringBuffer();
            for (String q : split) {
                long nowTime = System.currentTimeMillis();
                StringBuffer sb = new StringBuffer(URL);
                sb.append("q=");
                sb.append(URLEncoder.encode(q, "UTF-8"));
                sb.append("&from=auto");
                sb.append("&to=zh");
                sb.append("&");
                sb.append("appid=");
                sb.append(appid);
                sb.append("&");
                sb.append("salt=");
                sb.append(nowTime);
                sb.append("&");
                sb.append("sign=");
                StringBuffer saltBuffer = new StringBuffer();
                saltBuffer.append(appid);
                saltBuffer.append(q);
                saltBuffer.append(nowTime);
                saltBuffer.append(securityKey);
                String salt = CommonUtils.encryption(new String(saltBuffer));
                sb.append(salt);
                String jsonString = HTTP.getTransportResult(new String(sb));
                Map<String, Object> map = CommonUtils.strToMap(jsonString);
                if (!map.containsKey("error_code")) {
                    Object result = map.get("trans_result");
                    if (result instanceof List) {
                        List list = (List) result;
                        Map m = (Map) list.get(0);
                        if (!StringUtils.isEmpty(stringBuffer)) {
                            stringBuffer.append("\n");
                        }
                        stringBuffer.append(m.get("dst"));
                    }
                }
            }
            return new String(stringBuffer);
        } catch (Exception e) {
            return super.doTransport(text);
        }
    }
}
