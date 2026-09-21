import { Platform } from 'react-native';
import MainTabsIOS from './MainTabs.ios';
import MainTabsAndroid from './MainTabs.android';

const MainTabs = Platform.OS === 'ios' ? MainTabsIOS : MainTabsAndroid;

export default MainTabs;

