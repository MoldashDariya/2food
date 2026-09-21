import { Platform } from 'react-native';
import AdminTabsIOS from './AdminTabs.ios';
import AdminTabsAndroid from './AdminTabs.android';

const AdminTabs = Platform.OS === 'ios' ? AdminTabsIOS : AdminTabsAndroid;

export default AdminTabs;
