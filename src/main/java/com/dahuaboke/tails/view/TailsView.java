package com.dahuaboke.tails.view;

import javax.swing.*;

/**
 * author: dahua
 * date: 2023/12/13 16:10
 */
public class TailsView {

    private JTextField baiduAppIdText;
    private JTextField baiduSecretText;
    private JLabel baiduAppIdLabel;
    private JLabel baiduSecretLabel;
    private JPanel rootPanel;

    public JTextField getBaiduAppIdText() {
        return baiduAppIdText;
    }

    public JTextField getBaiduSecretText() {
        return baiduSecretText;
    }

    public JPanel getRootPanel() {
        return rootPanel;
    }
}
