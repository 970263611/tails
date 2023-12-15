package com.dahuaboke.tails;

import com.intellij.DynamicBundle;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.NotNull;
import org.jetbrains.annotations.PropertyKey;

import java.util.Arrays;
import java.util.function.Supplier;

/**
 * author: dahua
 * date: 2023/12/13 15:53
 */
public class TailsBundle extends DynamicBundle {

    public static final TailsBundle INSTANCE = new TailsBundle();

    private TailsBundle() {
        super("messages.TailsBundle");
    }

    public static final String message(@PropertyKey(resourceBundle = "messages.TailsBundle") @NotNull String key, @NotNull Object... params) {
        Intrinsics.checkNotNullParameter(key, "key");
        Intrinsics.checkNotNullParameter(params, "params");
        String var10000 = INSTANCE.getMessage(key, Arrays.copyOf(params, params.length));
        Intrinsics.checkNotNullExpressionValue(var10000, "getMessage(key, *params)");
        return var10000;
    }

    public static final Supplier messagePointer(@PropertyKey(resourceBundle = "messages.TailsBundle") @NotNull String key, @NotNull Object... params) {
        Intrinsics.checkNotNullParameter(key, "key");
        Intrinsics.checkNotNullParameter(params, "params");
        Supplier var10000 = INSTANCE.getLazyMessage(key, Arrays.copyOf(params, params.length));
        Intrinsics.checkNotNullExpressionValue(var10000, "getLazyMessage(key, *params)");
        return var10000;
    }
}
