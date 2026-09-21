import React, { useEffect, useMemo, useState } from 'react';
import {
    ScrollView,
    Text,
    View,
    StyleSheet,
    ActivityIndicator,
    TouchableOpacity,
    useWindowDimensions
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import {
    RestaurantCard,
    IconButton,
    FilterChip,
    FilterModal,
    Input
} from '../components';
import { useTranslation } from 'react-i18next';
import { colors, fontSizes, spacing } from '../constants/theme';
import { getAllRestaurants } from '../services/restaurants';
import { getProductsByRestaurant } from '../services/products';
import type { Restaurant, Product } from '../types/types';
import MultiSlider from '@ptomasroos/react-native-multi-slider';

type SortOption = 'default' | 'price' | 'time' | 'rating';
type TimeOption = 'all' | 'under20' | 'under30';

const SORT_KEYS: Record<SortOption, string> = {
    default: 'search.default',
    price: 'search.price',
    time: 'search.pickupTime',
    rating: 'search.rating',
};

const TIME_KEYS: Record<TimeOption, string> = {
    all: 'search.time',
    under20: 'search.under20',
    under30: 'search.under30',
};

const getReadyInMinutes = (pickupStart?: string) => {
    if (!pickupStart) return Number.MAX_SAFE_INTEGER;
    const pickupDate = new Date(pickupStart);
    if (Number.isNaN(pickupDate.getTime())) return Number.MAX_SAFE_INTEGER;
    return Math.max(0, Math.round((pickupDate.getTime() - Date.now()) / 60000));
};

export default function RestaurantsScreen() {
    const router = useRouter();
    const { t } = useTranslation();
    const { width } = useWindowDimensions();

    const [rawRestaurants, setRawRestaurants] = useState<Restaurant[]>([]);
    const [restaurantProducts, setRestaurantProducts] = useState<Record<string, Product[]>>({});
    const [loading, setLoading] = useState(true);

    const [sortOption, setSortOption] = useState<SortOption>('default');
    const [timeOption, setTimeOption] = useState<TimeOption>('all');
    
    const [minPrice, setMinPrice] = useState(0);
    const [maxPrice, setMaxPrice] = useState(10000);
    const [absoluteMax, setAbsoluteMax] = useState(10000);
    const [tempRange, setTempRange] = useState([0, 10000]);

    const [isSortModalVisible, setIsSortModalVisible] = useState(false);
    const [isPriceModalVisible, setIsPriceModalVisible] = useState(false);
    const [isTimeModalVisible, setIsTimeModalVisible] = useState(false);
    const [isMainFilterVisible, setIsMainFilterVisible] = useState(false);

    useEffect(() => {
        (async () => {
            try {
                const data = await getAllRestaurants();
                const rests = data.restaurants ?? [];
                setRawRestaurants(rests);

                const productsMap: Record<string, Product[]> = {};
                let globalMaxPrice = 0;

                await Promise.all(rests.map(async (r) => {
                    const res = await getProductsByRestaurant(r.id);
                    const activeProducts = (res.products ?? []).filter(p => p.is_active && p.quantity_available > 0);
                    productsMap[r.id] = activeProducts;
                    activeProducts.forEach(p => {
                        if (p.discounted_price > globalMaxPrice) globalMaxPrice = p.discounted_price;
                    });
                }));

                setRestaurantProducts(productsMap);
                const finalMax = globalMaxPrice || 10000;
                setAbsoluteMax(finalMax);
                setMaxPrice(finalMax);
                setTempRange([0, finalMax]);
            } catch (error) {
                console.error("Error loading restaurants:", error);
            } finally {
                setLoading(false);
            }
        })();
    }, []);

    const filteredData = useMemo(() => {
        let result = rawRestaurants.filter(rest => {
        const products = restaurantProducts[rest.id] ?? [];
            return products.some(p => {
                const matchesPrice = p.discounted_price >= minPrice && p.discounted_price <= maxPrice;
                const readyMin = getReadyInMinutes(p.pickup_time_start);
                const matchesTime = 
                timeOption === 'under20' ? readyMin <= 20 : 
                timeOption === 'under30' ? readyMin <= 30 : true;
                return matchesPrice && matchesTime;
            });
        });

        const sorted = [...result];
        if (sortOption === 'price') {
            sorted.sort((a, b) => {
                const minA = Math.min(...(restaurantProducts[a.id]?.map(p => p.discounted_price) ?? [0]));
                const minB = Math.min(...(restaurantProducts[b.id]?.map(p => p.discounted_price) ?? [0]));
                return minA - minB;
            });
        } else if (sortOption === 'time') {
            sorted.sort((a, b) => {
                const timeA = Math.min(...(restaurantProducts[a.id]?.map(p => getReadyInMinutes(p.pickup_time_start)) ?? [999]));
                const timeB = Math.min(...(restaurantProducts[b.id]?.map(p => getReadyInMinutes(p.pickup_time_start)) ?? [999]));
                return timeA - timeB;
            });
        }
        return sorted;
    }, [rawRestaurants, restaurantProducts, minPrice, maxPrice, timeOption, sortOption]);

    const activeFiltersCount = useMemo(() => {
        let count = 0;
        if (sortOption !== 'default') count++;
        if (timeOption !== 'all') count++;
        if (minPrice !== 0 || maxPrice !== absoluteMax) count++;
        return count;
    }, [sortOption, timeOption, minPrice, maxPrice, absoluteMax]);

    const handleReset = () => {
        setSortOption('default');
        setTimeOption('all');
        setMinPrice(0);
        setMaxPrice(absoluteMax);
        setTempRange([0, absoluteMax]);
    };

    const renderPriceContent = () => (
        <View>
            <View style={styles.priceInputsRow}>
                <View style={styles.priceInputWrapper}>
                    <Input
                        size="medium"
                        keyboardType="numeric"
                        value={String(tempRange[0])}
                        onChangeText={(text) => setTempRange([parseInt(text.replace(/[^0-9]/g, '')) || 0, tempRange[1]])}
                    />
                </View>
                <Text style={styles.priceDash}>-</Text>
                <View style={styles.priceInputWrapper}>
                    <Input
                        size="medium"
                        keyboardType="numeric"
                        value={String(tempRange[1])}
                        onChangeText={(text) => setTempRange([tempRange[0], parseInt(text.replace(/[^0-9]/g, '')) || absoluteMax])}
                    />
                </View>
            </View>
                <View style={styles.sliderWrapper}>
                    <MultiSlider
                    values={[tempRange[0], tempRange[1]]}
                    sliderLength={width - spacing.lg * 4}
                    onValuesChange={setTempRange}
                    min={0}
                    max={absoluteMax}
                    step={100}
                    selectedStyle={{ backgroundColor: colors.primary }}
                    markerStyle={styles.customMarker}
                    />
                </View>
        </View>
    );

    return (
        <View style={styles.container}>
        <SafeAreaView style={styles.safeArea} edges={['top']}>
            {/* Header */}
            <View style={styles.header}>
                <IconButton icon="chevron-back" onPress={() => router.back()} />
                <Text style={styles.headerTitle}>{t('restaurants.title')}</Text>
            </View>

            {/* Filters Bar */}
            <ScrollView horizontal showsHorizontalScrollIndicator={false} style={styles.filtersBar} contentContainerStyle={styles.filtersContent}>
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
                <FilterChip
                    label={minPrice !== 0 || maxPrice !== absoluteMax ? `${minPrice}-${maxPrice} ${t('common.currency')}` : t('search.price')}
                    onPress={() => setIsPriceModalVisible(true)}
                    active={minPrice !== 0 || maxPrice !== absoluteMax}
                />
                <FilterChip
                    label={t(TIME_KEYS[timeOption])}
                    onPress={() => setIsTimeModalVisible(true)}
                    active={timeOption !== 'all'}
                />
            </ScrollView>

            {loading ? (
                <ActivityIndicator size="large" color={colors.primary} style={styles.loader} />
                ) : (
                <ScrollView style={styles.scrollView} showsVerticalScrollIndicator={false} contentContainerStyle={styles.scrollContent}>
                    {filteredData.map((rest) => (
                    <RestaurantCard
                        key={rest.id}
                        name={rest.name}
                        imageUrl={rest.image_url}
                        variant="full"
                        onPress={() => router.push({
                        pathname: '/restaurant/[id]',
                        params: { id: rest.id, restaurant: JSON.stringify(rest) }
                        })}
                    />
                    ))}
                </ScrollView>
            )}

            {/* MODALS */}
            <FilterModal visible={isSortModalVisible} title={t('search.sortBy')} onClose={() => setIsSortModalVisible(false)}>
            {(Object.keys(SORT_KEYS) as SortOption[]).map((key) => (
                <TouchableOpacity key={key} style={styles.sortOptionRow} onPress={() => { setSortOption(key); setIsSortModalVisible(false); }}>
                    <View style={[styles.radioOuter, sortOption === key && styles.radioOuterActive]}>
                        {sortOption === key && <View style={styles.radioInner} />}
                    </View>
                    <Text style={styles.sortOptionText}>{t(SORT_KEYS[key])}</Text>
                </TouchableOpacity>
            ))}
            </FilterModal>

            <FilterModal
                visible={isPriceModalVisible}
                title={t('search.price')}
                onClose={() => setIsPriceModalVisible(false)}
                onApply={() => { setMinPrice(tempRange[0]); setMaxPrice(tempRange[1]); setIsPriceModalVisible(false); }}
                onReset={() => { handleReset(); setIsPriceModalVisible(false); }}
                >
                {renderPriceContent()}
            </FilterModal>

            <FilterModal visible={isTimeModalVisible} title={t('search.pickupTime')} onClose={() => setIsTimeModalVisible(false)} onApply={() => setIsTimeModalVisible(false)}>
                <View style={styles.chipsRow}>
                    {(Object.keys(TIME_KEYS) as TimeOption[]).map((key) => (
                        <FilterChip key={key} label={t(TIME_KEYS[key])} active={timeOption === key} onPress={() => setTimeOption(key)} />
                    ))}
                </View>
            </FilterModal>

            <FilterModal
                visible={isMainFilterVisible}
                title={t('search.filter')}
                onClose={() => setIsMainFilterVisible(false)}
                onApply={() => setIsMainFilterVisible(false)}
                onReset={handleReset}
            >
                <Text style={styles.filterSectionTitle}>{t('search.pickupTime')}</Text>
                <View style={[styles.chipsRow, { marginBottom: spacing.lg }]}>
                    {(Object.keys(TIME_KEYS) as TimeOption[]).map((key) => (
                        <FilterChip key={key} label={t(TIME_KEYS[key])} active={timeOption === key} onPress={() => setTimeOption(key)} />
                    ))}
                </View>
                <Text style={styles.filterSectionTitle}>{t('search.price')}</Text>
                {renderPriceContent()}
            </FilterModal>
        </SafeAreaView>
        </View>
    );
}

const styles = StyleSheet.create({
    container: { 
        flex: 1, 
        backgroundColor: colors.background 
    },
    safeArea: { 
        flex: 1 
    },
    header: { 
        flexDirection: 'row', 
        alignItems: 'center', 
        paddingHorizontal: spacing.lg, 
        paddingVertical: spacing.md 
    },
    headerTitle: { 
        fontSize: fontSizes.lg, 
        fontWeight: '600', 
        color: colors.black, 
        marginLeft: spacing.lg 
    },
    filtersBar: { 
        maxHeight: 50, 
    },
    filtersContent: { 
        paddingHorizontal: spacing.lg, 
        gap: spacing.sm, 
    },
    scrollView: { flex: 1 },
    scrollContent: { 
        paddingHorizontal: spacing.lg, 
        paddingBottom: 40, 
        gap: spacing.sm
    },
    loader: { 
        flex: 1, 
        justifyContent: 'center' 
    },
    priceInputsRow: { 
        flexDirection: 'row', 
        alignItems: 'center', 
        justifyContent: 'space-between', 
        marginBottom: spacing.lg 
    },
    priceInputWrapper: { 
        flex: 1 
    },
    priceDash: { 
        marginHorizontal: 12, 
        fontSize: 18 
    },
    sliderWrapper: { 
        alignItems: 'center', 
        justifyContent: 'center' 
    },
    customMarker: { 
        height: 24, 
        width: 24, 
        borderRadius: 12, 
        backgroundColor: colors.primary, 
        borderWidth: 2, 
        borderColor: colors.background 
    },
    sortOptionRow: { 
        flexDirection: 'row', 
        alignItems: 'center', 
        paddingVertical: spacing.lg, 
        gap: spacing.lg 
    },
    radioOuter: { 
        width: 24, 
        height: 24, 
        borderRadius: 12, 
        borderWidth: 2, 
        borderColor: colors.gray400, 
        justifyContent: 'center', 
        alignItems: 'center' 
    },
    radioOuterActive: { 
        borderColor: colors.primary 
    },
    radioInner: { 
        width: 12, 
        height: 12, 
        borderRadius: 6, 
        backgroundColor: colors.primary 
    },
    sortOptionText: { 
        fontSize: 16, 
        color: colors.black 
    },
    chipsRow: { 
        flexDirection: 'row', 
        flexWrap: 'wrap', 
        gap: spacing.sm 
    },
    filterSectionTitle: { 
        fontSize: fontSizes.sm, 
        fontWeight: '600', 
        color: colors.black, 
        marginBottom: spacing.sm 
    },
});