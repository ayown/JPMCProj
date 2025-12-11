-- Seed RBI circulars data
INSERT INTO rbi_circulars (id, circular_number, title, content, issued_date, effective_date, category, keywords, is_active, source_url, created_at, updated_at)
VALUES
    (
        uuid_generate_v4(),
        'RBI/2023-24/001',
        'KYC Compliance Guidelines for Banks',
        'All banks are required to conduct periodic KYC updates for existing customers. The KYC process must be completed within the specified timeline as per regulatory requirements.',
        '2023-04-01',
        '2023-05-01',
        'KYC',
        ARRAY['kyc', 'compliance', 'verification', 'customer'],
        TRUE,
        'https://www.rbi.org.in/Scripts/NotificationUser.aspx',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        uuid_generate_v4(),
        'RBI/2023-24/015',
        'Enhanced Security Measures for Digital Banking',
        'Banks must implement enhanced security measures including two-factor authentication, transaction alerts, and fraud detection systems for all digital banking channels.',
        '2023-06-15',
        '2023-07-01',
        'SECURITY',
        ARRAY['security', 'digital banking', 'authentication', 'fraud detection'],
        TRUE,
        'https://www.rbi.org.in/Scripts/NotificationUser.aspx',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        uuid_generate_v4(),
        'RBI/2023-24/028',
        'Customer Protection Guidelines',
        'Guidelines for protecting customers from fraudulent activities including phishing, vishing, and social engineering attacks. Banks must educate customers about safe banking practices.',
        '2023-09-10',
        '2023-10-01',
        'COMPLIANCE',
        ARRAY['customer protection', 'fraud', 'phishing', 'awareness'],
        TRUE,
        'https://www.rbi.org.in/Scripts/NotificationUser.aspx',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        uuid_generate_v4(),
        'RBI/2022-23/089',
        'Periodic KYC Update Requirements',
        'Customers must update their KYC details periodically. For high-risk customers, KYC must be updated annually. For medium-risk customers, once every two years.',
        '2022-12-01',
        '2023-01-01',
        'KYC',
        ARRAY['kyc', 'periodic update', 'customer verification'],
        TRUE,
        'https://www.rbi.org.in/Scripts/NotificationUser.aspx',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        uuid_generate_v4(),
        'RBI/2024-25/003',
        'SMS Banking Security Standards',
        'Banks must ensure that all SMS communications to customers include proper sender identification and security warnings. Customers should be advised never to share OTP or sensitive information.',
        '2024-01-15',
        '2024-02-01',
        'SECURITY',
        ARRAY['sms', 'security', 'otp', 'customer awareness'],
        TRUE,
        'https://www.rbi.org.in/Scripts/NotificationUser.aspx',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );

