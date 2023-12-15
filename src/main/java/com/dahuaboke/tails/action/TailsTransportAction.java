package com.dahuaboke.tails.action;

import com.dahuaboke.tails.TailsStateComponent;
import com.dahuaboke.tails.adapter.BaiduTransportAdapter;
import com.dahuaboke.tails.adapter.TransportAdapter;
import com.intellij.notification.*;
import com.intellij.openapi.actionSystem.AnAction;
import com.intellij.openapi.actionSystem.AnActionEvent;
import com.intellij.openapi.actionSystem.CommonDataKeys;
import com.intellij.openapi.application.ApplicationManager;
import com.intellij.openapi.command.WriteCommandAction;
import com.intellij.openapi.editor.Document;
import com.intellij.openapi.editor.Editor;
import com.intellij.openapi.editor.SelectionModel;
import org.apache.commons.lang3.ObjectUtils;
import org.apache.commons.lang3.StringUtils;
import org.jetbrains.annotations.NotNull;

/**
 * author: dahua
 * date: 2023/12/13 15:39
 */
public class TailsTransportAction extends AnAction {

    private static TailsStateComponent tailsStateComponent = ApplicationManager.getApplication().getService(TailsStateComponent.class);
    private static NotificationGroupManager notificationGroupManager = NotificationGroupManager.getInstance();
    private static NotificationGroup notificationGroup = notificationGroupManager.getNotificationGroup("transportNotification");

    static {
        String baiduAppId = tailsStateComponent.getBaiduAppId();
        String baiduSecret = tailsStateComponent.getBaiduSecret();
        if (StringUtils.isNoneEmpty(baiduAppId, baiduSecret)) {
            new BaiduTransportAdapter(baiduAppId, baiduSecret);
        }
    }

    private TransportAdapter transportAdapter = new TransportAdapter();

    @Override
    public void actionPerformed(@NotNull AnActionEvent e) {
        Editor editor = e.getRequiredData(CommonDataKeys.EDITOR);
        SelectionModel selectionModel = editor.getSelectionModel();
        String selectedText = selectionModel.getSelectedText();
        String result = transportAdapter.transport(selectedText);
        Document document = editor.getDocument();
        WriteCommandAction.runWriteCommandAction(e.getProject(), () -> {
            if (ObjectUtils.isNotEmpty(result)) {
                document.replaceString(selectionModel.getSelectionStart(), selectionModel.getSelectionEnd(), result);
            } else {
                Notification notification = notificationGroup.createNotification("翻译失败", NotificationType.ERROR);
                Notifications.Bus.notify(notification);
            }
        });
    }
}
