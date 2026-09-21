import React, { useEffect, useMemo, useState } from 'react';
import {
  FlatList,
  Text,
  View,
  StyleSheet,
  useWindowDimensions,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { FilterChip, ProductCard, SearchBar, RestaurantItem, EmptyState, FilterModal, Input } from '../../components';
import { borderRadius, colors, fontSizes, spacing } from '../../constants/theme';
import {
  applySearchFilters,
  type PriceOption,
  type SearchRestaurantEntity,
  type SortOption,
  type TimeOption,
} from '../../mocks/mockFlow';
import { getAllRestaurants } from '../../services/restaurants';
import { getProductsByRestaurant } from '../../services/products';
import type { Product, Restaurant } from '../../types/types';
import MultiSlider from '@ptomasroos/react-native-multi-slider';

const SORT_KEYS: Record<SortOption, string> = {
  default: 'search.default',
  price: 'search.price',
  time: 'search.pickupTime',
  rating: 'search.rating',
  discount: 'search.discount',
};

const PRICE_KEYS: Record<PriceOption, string> = {
  all: 'search.price',
  under2000: 'search.under2000',
  under3000: 'search.under3000',
};

const TIME_KEYS: Record<TimeOption, string> = {
  all: 'search.time',
  under20: 'search.under20',
  under30: 'search.under30',
};

const formatTime = (value?: string) => {
  if (!value) return '';
  const match = value.match(/(\d{2}:\d{2})/);
  return match ? match[1] : '';
};

const getReadyInMinutes = (pickupStart?: string) => {
  if (!pickupStart) return Number.MAX_SAFE_INTEGER;
  const pickupDate = new Date(pickupStart);
  if (Number.isNaN(pickupDate.getTime())) return Number.MAX_SAFE_INTEGER;
  return Math.max(0, Math.round((pickupDate.getTime() - Date.now()) / 60000));
};

export default function SearchScreen() {
  const router = useRouter();
  const { t } = useTranslation();
  const { width } = useWindowDimensions();
  const [searchText, setSearchText] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');
  const [sortOption, setSortOption] = useState<SortOption>('default');
  const [priceOption, setPriceOption] = useState<PriceOption>('all');
  const [timeOption, setTimeOption] = useState<TimeOption>('all');
  const [searchEntities, setSearchEntities] = useState<SearchRestaurantEntity[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const [isSortModalVisible, setIsSortModalVisible] = useState(false);

  const [isPriceModalVisible, setIsPriceModalVisible] = useState(false);
  const [isTimeModalVisible, setIsTimeModalVisible] = useState(false);

  const [tempPrice, setTempPrice] = useState<PriceOption>(priceOption);
  const [tempTime, setTempTime] = useState<TimeOption>(timeOption);

  const [isMainFilterVisible, setIsMainFilterVisible] = useState(false);

  const [priceRange, setPriceRange] = useState<[number, number] | null>(null);
  const [tempRange, setTempRange] = useState([0, 5000]);
  

  const { absoluteMinPrice, absoluteMaxPrice } = useMemo(() => {
    if (searchEntities.length === 0) return { absoluteMinPrice: 0, absoluteMaxPrice: 10000 };

    const allPrices = searchEntities.flatMap(e => 
      e.allProducts.map(p => p.discounted_price)
    );

    if (allPrices.length === 0) return { absoluteMinPrice: 0, absoluteMaxPrice: 10000 };

    return {
      absoluteMinPrice: Math.min(...allPrices),
      absoluteMaxPrice: Math.max(...allPrices)
    };
  }, [searchEntities]);
  

  const handleApplyPrice = () => {
    setPriceRange([tempRange[0], tempRange[1]]);
    setIsPriceModalVisible(false);
    setIsMainFilterVisible(false);
  };

  const handleSelectSort = (option: SortOption) => {
    setSortOption(option);
    setIsSortModalVisible(false);
  };

  useEffect(() => {
    if (isPriceModalVisible || isMainFilterVisible) {
      setTempRange(priceRange ?? [absoluteMinPrice, absoluteMaxPrice]);
    }
  }, [isPriceModalVisible, isMainFilterVisible]);

  useEffect(() => {
    const timeoutId = setTimeout(() => {
      setDebouncedSearch(searchText.trim().toLowerCase());
    }, 300);

    return () => clearTimeout(timeoutId);
  }, [searchText]);

  useEffect(() => {
  let isMounted = true;

  const loadRestaurants = async () => {
    setIsLoading(true);
    try {
      const data = await getAllRestaurants();
      const restaurants = data.restaurants ?? [];
      
      const productsResults = await Promise.allSettled(
        restaurants.map((restaurant) => getProductsByRestaurant(restaurant.id)),
      );

      const entities = restaurants.map((restaurant, index) => {
        const productsResult = productsResults[index];
        const products = productsResult.status === 'fulfilled'
          ? productsResult.value.products ?? []
          : [];

        const activeProducts = products.filter(p => p.is_active && p.quantity_available > 0);

        const offersCount = activeProducts.length;
        const totalSupplyLeft = activeProducts.reduce((sum, p) => sum + p.quantity_available, 0);

        return {
          id: restaurant.id,
          name: restaurant.name,
          address: restaurant.address ?? t('common.addressNotSpecified'),
          description: restaurant.description ?? t('common.noDescription'),
          logo: restaurant.image_url,
          image_url: restaurant.image_url,
          workingHours: '9:00-21:00',
          rating: restaurant.rating_avg ?? 0,
          allProducts: activeProducts, 
          filteredProducts: activeProducts,
          
          offersCount,      
          totalSupplyLeft,  
          
          bestOffer: activeProducts.length > 0 ? {
            id: activeProducts[0].id,
            title: activeProducts[0].name,
            image: activeProducts[0].image_url,
            originalPrice: activeProducts[0].original_price,
            discountedPrice: activeProducts[0].discounted_price,
            readyInMinutes: 20,
            supplyLeft: activeProducts[0].quantity_available,
          } : undefined
        };
      });

      if (isMounted) {
        setSearchEntities(entities);

        const allPrices = entities.flatMap(e => e.allProducts.map(p => p.discounted_price));
        if (allPrices.length > 0) {
          const min = Math.min(...allPrices);
          const max = Math.max(...allPrices);
          setTempRange([min, max]);
        }
      }
      } catch (error) {
        console.error("Ошибка загрузки данных поиска:", error);
        if (isMounted) setSearchEntities([]);
      } finally {
        if (isMounted) setIsLoading(false);
      }
    };

    loadRestaurants();
    return () => { isMounted = false; };
  }, []);

  const filteredData = useMemo(() => {
    const searchLower = debouncedSearch.toLowerCase();

    const result = searchEntities
      .map((restaurant) => {
        const matchedProducts = restaurant.allProducts.filter((product) => {
          const matchesSearch = 
            product.name.toLowerCase().includes(searchLower) || 
            restaurant.name.toLowerCase().includes(searchLower);

          const matchesPrice =
            priceRange === null ||
            (product.discounted_price >= priceRange[0] && product.discounted_price <= priceRange[1]);

          const readyMin = getReadyInMinutes(product.pickup_time_start);
          const matchesTime = 
            timeOption === 'under20' ? readyMin <= 20 : 
            timeOption === 'under30' ? readyMin <= 30 : true;

          return matchesSearch && matchesPrice && matchesTime;
        });

        return { ...restaurant, filteredProducts: matchedProducts };
      })
      .filter((res) => res.filteredProducts.length > 0);

    const sorted = [...result];

    if (sortOption === 'rating') {
      sorted.sort((a, b) => b.rating - a.rating);
    } 
    else if (sortOption === 'discount') {
      sorted.sort((a, b) => {
        const aBest = a.filteredProducts[0];
        const bBest = b.filteredProducts[0];
        const left = (aBest?.original_price ?? 0) - (aBest?.discounted_price ?? 0);
        const right = (bBest?.original_price ?? 0) - (bBest?.discounted_price ?? 0);
        return right - left;
      });
    }
    else if (sortOption === 'price') {
      sorted.sort((a, b) => 
        (a.filteredProducts[0]?.discounted_price ?? 0) - (b.filteredProducts[0]?.discounted_price ?? 0)
      );
    }
    else if (sortOption === 'time') {
      sorted.sort((a, b) => {
        const aTime = getReadyInMinutes(a.filteredProducts[0]?.pickup_time_start);
        const bTime = getReadyInMinutes(b.filteredProducts[0]?.pickup_time_start);
        return aTime - bTime;
      });
    }

    return sorted;
  }, [searchEntities, debouncedSearch, priceRange, timeOption, sortOption]);


  const activeFiltersCount = useMemo(() => {
    let count = 0;
    if (sortOption !== 'default') count++;
    if (timeOption !== 'all') count++;
    if (priceRange !== null) count++;
    return count;
  }, [sortOption, timeOption, priceRange]);

  const handleRestaurantPress = (restaurant: SearchRestaurantEntity) => {
    router.push({
      pathname: '/restaurant/[id]',
      params: {
        id: restaurant.id,
        restaurant: JSON.stringify(restaurant),
      }
    });
  };

  const handleProductPress = (
    product: Product,
    restaurant: SearchRestaurantEntity
  ) => {
    router.push({
      pathname: '/product/[id]',
      params: {
        id: product.id,
        product: JSON.stringify(product),
        restaurant: JSON.stringify(restaurant),
      },
    });
  };

  const handleReset = () => {
    setSortOption('default');
    setPriceOption('all');
    setTimeOption('all');
    setTempTime('all');
    setPriceRange(null);
    setTempRange([absoluteMinPrice, absoluteMaxPrice]);
  };

  const renderRestaurantSection = ({ item }: { item: SearchRestaurantEntity }) => (
    <RestaurantItem
      id={item.id}
      name={item.name}
      logo={item.logo}
      workingHours={item.workingHours}
      rating={item.rating}
      onPress={() => handleRestaurantPress(item)}
    >
      {/* Горизонтальный список отфильтрованных продуктов */}
      <FlatList
        data={item.filteredProducts}
        horizontal
        showsHorizontalScrollIndicator={false}
        keyExtractor={(product) => product.id}
        contentContainerStyle={styles.productsHorizontalList}
        snapToInterval={width * 0.7 + spacing.md} 
        decelerationRate="fast"
        renderItem={({ item: product }) => (
          <View style={styles.productCardWrap}>
            <ProductCard
              variant="simple"
              name={product.name}
              imageUrl={product.image_url}
              oldPrice={`${product.original_price} ₸`}
              price={`${product.discounted_price} ₸`}
              onPress={() => handleProductPress(product, item)}
            />
          </View>
        )}
      />
    </RestaurantItem>
  );

  const renderEmptyState = () => {
    if (isLoading) {
      return (
        <View style={{ flex: 1, alignItems: 'center', justifyContent: 'center', paddingTop: 80 }}>
          <ActivityIndicator size="large" color={colors.primary} />
        </View>
      );
    }
    if (searchText.length === 0) {
      return (
        <EmptyState
          title={t('search.emptyTitle')}
          description={t('search.emptyDescription')}
          imageSource={require('../../assets/pictures/search_empty.png')}
        />
      );
    }
    return (
      <EmptyState
        title={t('search.emptyTitle')}
        description={t('search.errorDescription')}
        imageSource={require('../../assets/pictures/search_empty.png')}
      />
    );
  }

  const renderPriceContent = () => {
    return (
      <View>
        {/* Инпуты */}
        <View style={styles.priceInputsRow}>
          <View style={styles.priceInputWrapper}>
            <Input
              size="medium"
              keyboardType="numeric"
              value={String(tempRange[0])}
              onChangeText={(text) => {
                const val = parseInt(text.replace(/[^0-9]/g, '')) || 0;
                setTempRange([val, tempRange[1]]);
              }}
              showCancelButton={false}
            />
          </View>

          <Text style={styles.priceDash}>-</Text>

          <View style={styles.priceInputWrapper}>
            <Input
              size="medium"
              keyboardType="numeric"
              value={String(tempRange[1])}
              onChangeText={(text) => {
                const val = parseInt(text.replace(/[^0-9]/g, '')) || absoluteMaxPrice;
                setTempRange([tempRange[0], val]);
              }}
              showCancelButton={false}
            />
          </View>
        </View>

        {/* Слайдер с динамическими границами */}
        <View style={styles.sliderWrapper}>
          <MultiSlider
            values={[tempRange[0], tempRange[1]]}
            sliderLength={width - spacing.lg * 4}
            onValuesChange={(values) => {
              setTempRange(values);
            }}
            min={absoluteMinPrice} 
            max={absoluteMaxPrice} 
            step={100}
            allowOverlap={false}
            snapped
            trackStyle={styles.sliderTrack}
            selectedStyle={styles.sliderSelectedTrack}
            unselectedStyle={styles.sliderUnselectedTrack}
            markerStyle={styles.customMarker}
          />
        </View>
      </View>
    );
  };

  const priceChipLabel = useMemo(() => {
    if (priceRange !== null) {
      return t('search.priceWithRange', { min: priceRange[0], max: priceRange[1] });
    }
    return t('search.price');
  }, [priceRange, t]);

  return (
    <View style={styles.container}>
      <SafeAreaView style={styles.safeArea} edges={['top']}>
        <View style={styles.searchContainer}>
          <SearchBar
            value={searchText}
            onChangeText={setSearchText}
            placeholder={t('search.placeholder')}
          />
        </View>

        <View style={styles.filtersContainer}>
          <FilterChip
            iconName="swap-vertical"
            iconOnly
            onPress={() => setIsSortModalVisible(true)}
            active={sortOption !== 'default'}
          />
          <FilterChip
            iconName="options-outline"
            iconOnly
            onPress={() => setIsMainFilterVisible(true)}
            active={activeFiltersCount > 0}
            hasNotification={activeFiltersCount > 0}
          />

          {/* Цена */}
          <FilterChip 
            label={priceChipLabel} 
            onPress={() => {
              setTempRange(priceRange ?? [absoluteMinPrice, absoluteMaxPrice]);
              setIsPriceModalVisible(true);
            }} 
           
          active={priceRange !== null}
          />

          {/* Время */}
          <FilterChip
            label={t(TIME_KEYS[timeOption])}
            onPress={() => setIsTimeModalVisible(true)}
            active={timeOption !== 'all'}
  />
        </View>

        <FlatList
          data={filteredData}
          keyExtractor={(item) => item.id}
          renderItem={renderRestaurantSection}
          ListEmptyComponent={renderEmptyState}
          contentContainerStyle={styles.listContent}
          showsVerticalScrollIndicator={false}
          keyboardShouldPersistTaps="handled"
        />

        <FilterModal
          visible={isSortModalVisible}
          title={t('search.sortBy')}
          onClose={() => setIsSortModalVisible(false)}
        >
          <View style={styles.sortOptionsContainer}>
            {(Object.keys(SORT_KEYS) as SortOption[]).map((key) => {
              const isSelected = sortOption === key;
              const label = t(SORT_KEYS[key]);
              
              return (
                <TouchableOpacity
                  key={key}
                  style={styles.sortOptionRow}
                  onPress={() => handleSelectSort(key)}
                >
                  <View style={[styles.radioOuter, isSelected && styles.radioOuterActive]}>
                    {isSelected && <View style={styles.radioInner} />}
                  </View>
                  <Text style={[styles.sortOptionText, isSelected && styles.sortOptionTextActive]}>
                    {label}
                  </Text>
                </TouchableOpacity>
              );
            })}
          </View>
        </FilterModal>

        <FilterModal
          visible={isPriceModalVisible}
          title={t('search.price')}
          onClose={() => setIsPriceModalVisible(false)}
          onReset={() => {
            setPriceRange(null);
            setTempRange([absoluteMinPrice, absoluteMaxPrice]);
            setIsPriceModalVisible(false);
          }}
          onApply={handleApplyPrice}
        >
          {renderPriceContent()}
        </FilterModal>


        <FilterModal
          visible={isTimeModalVisible}
          title={t('search.pickupTime')}
          onClose={() => setIsTimeModalVisible(false)}
          onReset={() => {
            setTimeOption('all');
            setIsTimeModalVisible(false);
          }}
          onApply={() => setIsTimeModalVisible(false)}
        >
          <View style={styles.mainFilterContent}>
            <View style={styles.chipsRow}>
              {(Object.keys(TIME_KEYS) as TimeOption[]).map((key) => (
                <FilterChip
                  key={key}
                  label={t(TIME_KEYS[key])}
                  active={timeOption === key}
                  onPress={() => {
                    setTimeOption(prev => (prev === key && key !== 'all' ? 'all' : key));
                  }}
                />
              ))}
            </View>
          </View>
        </FilterModal>

        <FilterModal
          visible={isMainFilterVisible}
          title={t('search.filter')}
          onClose={() => setIsMainFilterVisible(false)}
          onApply={() => setIsMainFilterVisible(false)}
          onReset={handleReset}
        >
          <View style={styles.mainFilterContent}>
            {/* Секция Времени */}
            <Text style={styles.filterSectionTitle}>{t('search.pickupTime')}</Text>
            <View style={styles.chipsRow}>
              {(Object.keys(TIME_KEYS) as TimeOption[]).map((key) => (
                <FilterChip
                  key={key}
                  label={t(TIME_KEYS[key])}
                  active={timeOption === key}
                  onPress={() => {
                    setTimeOption(prev => (prev === key && key !== 'all' ? 'all' : key));
                  }}
                />
              ))}
            </View>

            {/* Секция Цены */}
            <Text style={[styles.filterSectionTitle]}>{t('search.price')}</Text>
            {renderPriceContent()}
          </View>
        </FilterModal>
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f4f4f4',
  },
  safeArea: {
    flex: 1,
  },
  listContent: {
    paddingBottom: 120,
    paddingHorizontal: spacing.lg,
  },
  searchContainer: {
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
  },
  filtersContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
    paddingHorizontal: spacing.lg,
    marginTop: spacing.lg,
    marginBottom: spacing.lg,
  },
  productCardWrap: {
    borderRadius: borderRadius.lg,
    alignSelf: 'flex-start',
  },

  sortOptionsContainer: {
    paddingBottom: spacing.lg,
  },
  sortOptionRow: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: spacing.lg,
    gap: spacing.lg,
  },
  radioOuter: {
    width: 24,
    height: 24,
    borderRadius: 14,
    borderWidth: 2,
    borderColor: colors.gray400,
    justifyContent: 'center',
    alignItems: 'center',
  },
  radioOuterActive: {
    borderColor: colors.primary, 
  },
  radioInner: {
    width: 12,
    height: 12,
    borderRadius: 6,
    backgroundColor: colors.primary,
  },
  sortOptionText: {
    fontSize: 16,
    color: colors.black,
    fontWeight: '400',
  },
  sortOptionTextActive: {
    fontWeight: '500',
  },
  priceInputsRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: spacing.lg,
  },
  priceInputWrapper: {
    flex: 1,
  },
  mainFilterContent: {
    gap: spacing.xxl
  },
  filterSectionTitle: {
    fontSize: fontSizes.sm,
    fontWeight: '600',
    color: colors.black,
  },
  chipsRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: spacing.sm,
  },
  priceDash: {
    marginHorizontal: 12,
    fontSize: 18,
    color: colors.black,
  },
  sliderWrapper: {
    alignItems: 'center',
    justifyContent: 'center',
  },
  sliderTrack: {
    height: 4,
    borderRadius: 2,
  },
  sliderSelectedTrack: {
    backgroundColor: colors.primary,
  },
  sliderUnselectedTrack: {
    backgroundColor: colors.gray400,
  },
  customMarker: {
    height: 24,
    width: 24,
    borderRadius: 12,
    backgroundColor: colors.primary, 
    shadowColor: 'transparent',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0,
    shadowRadius: 0,
    elevation: 0, 
    borderWidth: 1,
    borderColor: colors.background,
  },
  productsHorizontalList: {
    paddingLeft: 0,
    paddingRight: spacing.lg,
    gap: spacing.md, 
    marginTop: spacing.sm,
  },
});
