package com.dahuaboke.tails;

import com.dahuaboke.tails.adapter.BaiduTransportAdapter;
import com.dahuaboke.tails.adapter.TransportAdapter;
import com.dahuaboke.tails.view.TailsView;
import com.intellij.openapi.application.ApplicationManager;
import com.intellij.openapi.options.SearchableConfigurable;
import com.intellij.openapi.util.NlsContexts;
import org.apache.commons.lang3.StringUtils;
import org.jetbrains.annotations.NonNls;
import org.jetbrains.annotations.NotNull;

import javax.swing.*;

/**
 * author: dahua
 * date: 2023/12/13 17:48
 */
public class TailConfiguration implements SearchableConfigurable {

    private TailsStateComponent tailsStateComponent = ApplicationManager.getApplication().getService(TailsStateComponent.class);
    private TailsView view = new TailsView();

    @Override
    public @NotNull @NonNls String getId() {
        return "tailsView";
    }

    @Override
    public @NlsContexts.ConfigurableName String getDisplayName() {
        return "tails";
    }

    @Override
    public JComponent createComponent() {
        String baiduAppId = tailsStateComponent.getBaiduAppId();
        String baiduSecret = tailsStateComponent.getBaiduSecret();
        JTextField baiduAppIdText = view.getBaiduAppIdText();
        JTextField baiduSecretText = view.getBaiduSecretText();
        baiduAppIdText.setText(baiduAppId);
        baiduSecretText.setText(baiduSecret);
        return view.getRootPanel();
    }

    @Override
    public boolean isModified() {
        TailsStateComponent state = tailsStateComponent.getState();
        String baiduAppId = state.getBaiduAppId();
        String baiduAppIdText = view.getBaiduAppIdText().getText();
        if (!StringUtils.equals(baiduAppId, baiduAppIdText)) {
            return true;
        }
        String baiduSecret = state.getBaiduSecret();
        String baiduSecretText = view.getBaiduSecretText().getText();
        if (!StringUtils.equals(baiduSecret, baiduSecretText)) {
            return true;
        }
        return false;
    }

    @Override
    public void apply() {
        tailsStateComponent.setBaiduAppId(view.getBaiduAppIdText().getText());
        tailsStateComponent.setBaiduSecret(view.getBaiduSecretText().getText());
    }
}
