package com.dahuaboke.tails.util;

import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;

import java.io.IOException;

/**
 * author: dahua
 * date: 2023/12/13 15:42
 */
public class HTTP {

    private static final OkHttpClient okHttpClient = new OkHttpClient();
    private static final Request.Builder BUILDER = new Request.Builder();

    public static String getTransportResult(String url) throws Exception {
        Request request = BUILDER.get().url(url.toString()).build();
        Response response = okHttpClient.newCall(request).execute();
        if (response.code() == 200) {
            return response.body().string();
        }
        throw new Exception();
    }
}
