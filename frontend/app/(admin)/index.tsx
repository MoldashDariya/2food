import React, { useMemo, useEffect, useState, useCallback, useRef } from 'react';
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  TouchableOpacity,
  Alert,
  RefreshControl,
  ScrollView,
  useWindowDimensions,
} from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { useFocusEffect, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, fonts, fontSizes, spacing, borderRadius } from '../../constants/theme';
import { ADMIN_COLORS, ADMIN_STRINGS } from '../../constants/admin';
import { ProductListItem, ProductActionSheet } from '../../components/admin';
import { useAdminStore } from '../../stores/adminStore';
import {
  deletePartnerProduct,
  getPartnerCategories,
  getPartnerProducts,
} from '../../services/admin';
import type { Product } from '../../types/types';
import { IconButton, EmptyState, SearchBar, LoadingScreen } from '../../components';
import { getLocalizedCategoryLabel } from '../../utils/categories';

interface ProductViewModel {
  id: string;
  name: string;
  categoryId: string | number;
  categoryName?: string;
  categoryIcon?: string;
  imageUrl?: string;
  originalPrice: number;
  discountedPrice: number;
  quantityAvailable: number;
  isActive: boolean;
}

type TabType = 'active' | 'inactive';

export default function ProductsScreen() {
  const { t } = useTranslation();
  const insets = useSafeAreaInsets();
  const router = useRouter();
  const { width: screenWidth } = useWindowDimensions();
  const pagerRef = useRef<ScrollView>(null);
  const { selectedRestaurantId, syncRestaurants, restaurants } = useAdminStore();

  const [products, setProducts] = useState<ProductViewModel[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string | number>('all');
  const [refreshing, setRefreshing] = useState(false);
  const [loading, setLoading] = useState(true);
  const [searchText, setSearchText] = useState('');
  const [activeTab, setActiveTab] = useState<TabType>('active');
  const [selectedProduct, setSelectedProduct] = useState<ProductViewModel | null>(null);
  const [isMenuVisible, setMenuVisible] = useState(false);

  const loadData = useCallback(async () => {
    if (!selectedRestaurantId) return;
    
    let productsList: Product[] = [];
    let categoriesList: any[] = [];

    try {
      const [pResult, cResult] = await Promise.allSettled([
        getPartnerProducts(selectedRestaurantId),
        getPartnerCategories(),
      ]);

      if (pResult.status === 'fulfilled') {
        productsList = pResult.value.products || [];
      } else {
        console.error('Failed to load products:', pResult.reason);
      }

      if (cResult.status === 'fulfilled') {
        categoriesList = cResult.value.product_categories || cResult.value.categories || [];
      } else {
        console.error('Failed to load categories:', cResult.reason);
      }
    } catch (error) {
      console.error('Error in loadData:', error);
    }

    const categoryById = new Map<string | number, { name: string; icon: string }>();
    
    categoriesList.forEach((c) => {
      categoryById.set(c.id, {
        name: getLocalizedCategoryLabel({ name: c.name, label: c.label }),
        icon: c.icon_url ? '🏷️' : '📦',
      });
    });

    // "Hack": Discover categories from products if they are missing from categoriesList
    productsList.forEach((p) => {
      if (p.category_id !== undefined && p.category_id !== null && !categoryById.has(p.category_id)) {
        categoryById.set(p.category_id, {
          name: `${t('admin.product.category')} ${p.category_id}`,
          icon: '📦',
        });
      }
    });

    const mapped = productsList.map((p: Product) => {
      const catId = p.category_id ?? 'uncategorized';
      const catInfo = categoryById.get(catId);
      
      return {
        id: p.id,
        name: p.name,
        categoryId: catId,
        categoryName: catInfo?.name || t('admin.extras.noCategory'),
        categoryIcon: catInfo?.icon || '📦',
        imageUrl: p.image_url,
        originalPrice: p.original_price,
        discountedPrice: p.discounted_price,
        quantityAvailable: p.quantity_available,
        isActive: p.is_active,
      };
    });

    setProducts(mapped);
  }, [selectedRestaurantId]);

  useEffect(() => {
    syncRestaurants();
  }, [syncRestaurants]);

  useEffect(() => {
    let mounted = true;
    const run = async () => {
      if (!selectedRestaurantId) {
        setLoading(false);
        return;
      }
      setLoading(true);
      try {
        await loadData();
      } finally {
        if (mounted) setLoading(false);
      }
    };
    run();
    return () => {
      mounted = false;
    };
  }, [selectedRestaurantId, loadData]);

  useFocusEffect(
    useCallback(() => {
      if (selectedRestaurantId) {
        loadData();
      }
    }, [selectedRestaurantId, loadData]),
  );

  const categoryOptions = useMemo(() => {
    const base = [{ id: 'all', name: t('admin.extras.all'), icon: '📋' }];

    const productsInCurrentTab = products.filter((p) => {
      const isProductActive = p.isActive && p.quantityAvailable > 0;
      return activeTab === 'active' ? isProductActive : !isProductActive;
    });

    const unique = new Map<string | number, { id: string | number; name: string; icon: string }>();
    
    for (const p of productsInCurrentTab) {
      if (!unique.has(p.categoryId)) {
        unique.set(p.categoryId, {
          id: p.categoryId,
          name: p.categoryName || t('admin.extras.noCategory'),
          icon: p.categoryIcon || '📦',
        });
      }
    }

    return [...base, ...Array.from(unique.values())];
  }, [products, activeTab]);

  useEffect(() => {
    const isCategoryStillAvailable = categoryOptions.some(cat => cat.id === selectedCategory);
    if (!isCategoryStillAvailable) {
      setSelectedCategory('all');
    }
  }, [categoryOptions, selectedCategory]);

  const activeProducts = useMemo(() => {
    return products.filter((p) => {
      const matchesCategory = selectedCategory === 'all' || p.categoryId === selectedCategory;
      const matchesSearch = p.name.toLowerCase().includes(searchText.toLowerCase().trim());
      const isProductActive = p.isActive && p.quantityAvailable > 0;

      return matchesCategory && matchesSearch && isProductActive;
    });
  }, [products, selectedCategory, searchText]);

  const inactiveProducts = useMemo(() => {
    return products.filter((p) => {
      const matchesCategory = selectedCategory === 'all' || p.categoryId === selectedCategory;
      const matchesSearch = p.name.toLowerCase().includes(searchText.toLowerCase().trim());
      const isProductActive = p.isActive && p.quantityAvailable > 0;

      return matchesCategory && matchesSearch && !isProductActive;
    });
  }, [products, selectedCategory, searchText]);

  const handleRefresh = useCallback(async () => {
    if (!selectedRestaurantId) return;
    setRefreshing(true);
    try {
      await loadData();
    } finally {
      setRefreshing(false);
    }
  }, [selectedRestaurantId, loadData]);

  const handleOpenMenu = (product: ProductViewModel) => {
    setSelectedProduct(product);
    setMenuVisible(true);
  };

  const handleStatusChange = async () => {
    if (!selectedProduct || !selectedRestaurantId) return;
    await loadData();
  };

  const handleEdit = (product: ProductViewModel) => {
    router.push({
      pathname: '/(admin)/add-product',
      params: { productId: product.id },
    });
  };

  const handleDelete = (product: ProductViewModel) => {
    if (!selectedRestaurantId) return;

    Alert.alert(
      ADMIN_STRINGS.product.delete,
      `${ADMIN_STRINGS.common.confirm} "${product.name}"?`,
      [
        { text: ADMIN_STRINGS.common.cancel, style: 'cancel' },
        {
          text: ADMIN_STRINGS.product.delete,
          style: 'destructive',
          onPress: async () => {
            try {
              await deletePartnerProduct(selectedRestaurantId, product.id);
              setProducts((prev) => prev.filter((p) => p.id !== product.id));
            } catch {
              Alert.alert(t('common.error'), t('admin.extras.deleteProductError'));
            }
          },
        },
      ],
    );
  };

  const renderCategoryFilter = () => (
    <View style={styles.filterContainer}>
    <FlatList
      horizontal
      showsHorizontalScrollIndicator={false}
      data={categoryOptions}
      keyExtractor={(item) => String(item.id)}
      contentContainerStyle={styles.filterList}
      renderItem={({ item }) => (
        <TouchableOpacity
          style={[
            styles.filterChip,
            selectedCategory === item.id && styles.filterChipActive,
          ]}
          onPress={() => setSelectedCategory(item.id)}
          activeOpacity={0.9}
        >
          <Text
            style={[
              styles.filterText,
              selectedCategory === item.id && styles.filterTextActive,
            ]}
          >
            {item.name}
          </Text>
        </TouchableOpacity>
      )}
    />
  </View>
  );

  const renderEmptyState = (statusLabel: string, currentListLength: number) => {
    if (products.length === 0) {
      return (
        <View style={styles.emptyContainer}>
          <EmptyState
            title={t('admin.extras.emptyMenuTitle')}
            description={t('admin.extras.emptyMenuDescription')}
            imageSource={require('../../assets/pictures/menu_empty.png')}
            buttonTitle={t('admin.extras.addProduct')}
            onButtonPress={() => router.push('/(admin)/add-product')}
          />
        </View>
      );
    }

    if (searchText.trim().length > 0 && currentListLength === 0) {
      return (
        <View style={styles.emptyContainer}>
          <EmptyState
            title={t('admin.extras.nothingFound')}
            description={t('search.errorDescription')}
            imageSource={require('../../assets/pictures/search_empty.png')}
          />
        </View>
      );
    }
    
    return (
      <View style={styles.emptyContainer}>
        <EmptyState
          title={t('admin.extras.categoryEmpty')}
          description={t('restaurant.emptyCategoryDescription')}
          imageSource={require('../../assets/pictures/menu_empty.png')}
        />
      </View>
    );
  };

  const currentRestaurant = restaurants.find((r) => r.id === selectedRestaurantId);

  return (
    <View style={[styles.container, { paddingTop: insets.top }]}>
      <View style={styles.header}>
        <View>
          <Text style={styles.title}>{ADMIN_STRINGS.nav.menu}</Text>
          {currentRestaurant ? (
            <Text style={styles.subtitle}>{currentRestaurant.name} {currentRestaurant.address ? `, ${currentRestaurant.address}` : ''}</Text>
          ) : null}
        </View>
        <IconButton
            icon="notifications"
            iconSize={18}
        />
      </View>

      <View style={styles.searchContainer}>
        <SearchBar
          value={searchText}
          onChangeText={setSearchText}
          placeholder={t('search.placeholder')}
        />
      </View>
    
      <View style={styles.tabContainer}>
        <TouchableOpacity 
          style={[styles.tab, activeTab === 'active' && styles.activeTab]} 
          onPress={() => {
            setActiveTab('active');
            pagerRef.current?.scrollTo({ x: 0, animated: true });
          }}
          activeOpacity={0.9}
        >
          <Text style={[styles.tabText, activeTab === 'active' && styles.activeTabText]}>
            Активные
          </Text>
        </TouchableOpacity>
        
        <TouchableOpacity 
          style={[styles.tab, activeTab === 'inactive' && styles.activeTab]} 
          onPress={() => {
            setActiveTab('inactive');
            pagerRef.current?.scrollTo({ x: screenWidth, animated: true });
          }}
          activeOpacity={0.9}
        >
          <Text style={[styles.tabText, activeTab === 'inactive' && styles.activeTabText]}>
            Неактивные
          </Text>
        </TouchableOpacity>
      </View>

      {renderCategoryFilter()}

      {loading ? (
        <LoadingScreen />
      ) : (
        <ScrollView
          ref={pagerRef}
          horizontal
          pagingEnabled
          showsHorizontalScrollIndicator={false}
          scrollEventThrottle={16}
          style={{ flex: 1 }}
          onMomentumScrollEnd={(e) => {
            const pageIndex = Math.round(e.nativeEvent.contentOffset.x / screenWidth);
            setActiveTab(pageIndex === 0 ? 'active' : 'inactive');
          }}
        >
          <View style={{ width: screenWidth, flex: 1 }}>
            <FlatList
              data={activeProducts}
              keyExtractor={(item) => item.id}
              contentContainerStyle={styles.listContent}
              refreshControl={
                <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} tintColor={ADMIN_COLORS.primary} />
              }
              ListEmptyComponent={() => renderEmptyState('active', activeProducts.length)}
              renderItem={({ item }) => (
                <ProductListItem
                  id={item.id}
                  name={item.name}
                  imageUrl={item.imageUrl}
                  originalPrice={item.originalPrice}
                  discountedPrice={item.discountedPrice}
                  quantity={item.quantityAvailable}
                  isActive={item.isActive}
                  onEdit={() => handleOpenMenu(item)}
                />
              )}
            />
          </View>
          <View style={{ width: screenWidth, flex: 1 }}>
            <FlatList
              data={inactiveProducts}
              keyExtractor={(item) => item.id}
              contentContainerStyle={styles.listContent}
              refreshControl={
                <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} tintColor={ADMIN_COLORS.primary} />
              }
              ListEmptyComponent={() => renderEmptyState('inactive', inactiveProducts.length)}
              renderItem={({ item }) => (
                <ProductListItem
                  id={item.id}
                  name={item.name}
                  imageUrl={item.imageUrl}
                  originalPrice={item.originalPrice}
                  discountedPrice={item.discountedPrice}
                  quantity={item.quantityAvailable}
                  isActive={item.isActive}
                  onEdit={() => handleOpenMenu(item)}
                />
              )}
            />
          </View>
        </ScrollView>
      )}

      <ProductActionSheet
        isVisible={isMenuVisible}
        onClose={() => setMenuVisible(false)}
        onEdit={() => selectedProduct && handleEdit(selectedProduct)}
        onChangeStatus={handleStatusChange}
        onDelete={() => selectedProduct && handleDelete(selectedProduct)}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.sm,
    paddingBottom: spacing.sm,
  },
  title: {
    fontSize: fontSizes.lg,
    fontFamily: fonts.semiBold,
    color: colors.black,
  },
  subtitle: {
    marginTop: 1,
    fontSize: fontSizes.sm,
    fontFamily: fonts.medium,
    color: colors.gray600,
  },
  tabContainer: {
    flexDirection: 'row',
    borderBottomWidth: 1,
    borderBottomColor: colors.gray200,
    marginBottom: spacing.lg,
    paddingHorizontal: spacing.lg,
  },
  tab: {
    alignItems: 'center',
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.xs,
    marginRight: spacing.xl,
    borderBottomWidth: 2,
    borderBottomColor: 'transparent',
  },
  activeTab: {
    borderBottomColor: colors.primary,
  },
  tabText: {
    fontSize: fontSizes.md,
    fontFamily: fonts.medium,
    color: colors.gray600,
  },
  activeTabText: {
    color: colors.primary,
    fontFamily: fonts.semiBold,
  },
  filterContainer: {
    marginBottom: spacing.lg,
  },
  filterList: {
    paddingHorizontal: spacing.lg,
    gap: spacing.sm, 
    flexDirection: 'row',
  },
  filterChip: {
    paddingHorizontal: 12,
    paddingVertical: 8,
    borderRadius: 12, 
    borderWidth: 1,
    borderColor: colors.gray400,
    backgroundColor: colors.background,
  },
  filterChipActive: {
    backgroundColor: '#F6FFF3', 
    borderColor: colors.primary, 
  },
  filterText: {
    fontSize: fontSizes.md,
    fontWeight: '500',
    color: colors.black,
  },
  filterTextActive: {
    color: colors.primary,
    fontWeight: '600',
  },
  listContent: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.xxl,
  },
  emptyContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingTop: 100,
  },
  emptyText: {
    fontSize: fontSizes.lg,
    fontFamily: fonts.semiBold,
    color: colors.gray600,
    marginTop: spacing.md,
  },
  searchContainer: {
    paddingHorizontal: spacing.lg,
    marginBottom: spacing.lg,
    marginTop: 8,
  },
});
