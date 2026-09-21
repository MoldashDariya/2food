import React from 'react';
import AdminTabs from '../../components/admin/AdminTabs';

// Admin (restaurant staff) screens share one tab bar.
// Metro picks AdminTabs.ios.tsx / AdminTabs.android.tsx by platform
export default function AdminLayout() {
    return <AdminTabs />;
}
