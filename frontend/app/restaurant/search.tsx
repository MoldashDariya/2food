import React, { useState, useMemo, useEffect } from 'react';
import { ScrollView, Text, View, StyleSheet } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { IconButton, SearchBar, ProductCard, EmptyState, LoadingScreen } from '../../components';
import { colors, fontSizes, spacing } from '../../constants/theme';
import { getProductsByRestaurant } from '../../services/products';
import type { Product } from '../../types/types';

export default function SearchScreen() {
    const router = useRouter();
    const { t } = useTranslation();
    const { id } = useLocalSearchParams<{ id: string }>();

    const [searchText, setSearchText] = useState('');
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!id) {
            console.warn("SearchScreen: Restaurant ID is missing in params");
            setLoading(false);
            return;
        }

        const loadData = async () => {
            try {
                setLoading(true);
                const result = await getProductsByRestaurant(id);
                const activeProducts = (result.products || []).filter(
                    p => p.quantity_available > 0 && p.is_active !== false
                );
                setProducts(activeProducts);
            } catch (error) {
                console.error("Search fetch error:", error);
            } finally {
                setLoading(false);
            }
        };

        loadData();
    }, [id]);

    const filteredProducts = useMemo(() => {
        const query = searchText.trim().toLowerCase();
        if (!query) return [];
        return products.filter(p => 
            p.name.toLowerCase().includes(query) || 
            (p.description && p.description.toLowerCase().includes(query))
        );
    }, [searchText, products]);

    const handleProductPress = (product: Product) => {
        router.push({
            pathname: '/product/[id]',
            params: { 
                id: product.id, 
                productString: JSON.stringify(product) 
            }
        });
    };

    if (loading) {
        return <LoadingScreen />;
    }

    return (
        <View style={styles.container}>
            <SafeAreaView style={styles.safeArea} edges={['top']}>
                {/* Header */}
                <View style={styles.header}>
                    <IconButton icon="chevron-back" onPress={() => router.back()} />
                    <Text style={styles.headerTitle}>{t('search.title')}</Text>
                </View>

                <View style={styles.searchContainer}>
                    <SearchBar
                        value={searchText}
                        onChangeText={setSearchText}
                        placeholder={t('search.placeholder')}
                        size='small'
                    />
                </View>

                <ScrollView
                    style={styles.scrollView}
                    showsVerticalScrollIndicator={false}
                    contentContainerStyle={styles.scrollContent}
                >
                    {searchText.length > 0 && filteredProducts.length > 0 && (
                        <Text style={styles.resultsTitle}>{t('search.resultsTitle')}</Text>
                    )}

                    {filteredProducts.map((product) => (
                        <ProductCard
                            key={product.id}
                            variant='compact'
                            name={product.name}
                            details={product.description || t('product.fallbackDetails')}
                            oldPrice={`${product.original_price} ${t('common.currency')}`}
                            price={`${product.discounted_price} ${t('common.currency')}`}
                            imageUrl={product.image_url}
                            remainingCount={product.quantity_available}
                            onPress={() => handleProductPress(product)}
                        />
                    ))}

                    {searchText.length > 0 && filteredProducts.length === 0 && (
                        <View style={styles.emptyBox}>
                            <EmptyState
                                title={t('search.emptyTitle')}
                                description={t('search.errorDescription')}
                                imageSource={require('../../assets/pictures/search_empty.png')}
                            />
                        </View>
                    )}

                    {searchText.length === 0 && (
                        <EmptyState
                            title={t('search.enterNameTitle')}
                            description={t('search.enterNameDescription')}
                            imageSource={require('../../assets/pictures/search_empty.png')}
                        />
                    )}
                </ScrollView>
            </SafeAreaView>
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: colors.background,
    },
    safeArea: {
        flex: 1,
    },
    header: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingHorizontal: spacing.lg,
        paddingTop: spacing.lg,
        marginBottom: spacing.lg,
    },
    headerTitle: {
        fontSize: fontSizes.lg,
        fontWeight: '600',
        color: colors.black,
        marginLeft: spacing.lg,
    },
    searchContainer: {
        paddingHorizontal: spacing.lg,
    },
    resultsTitle: {
        fontSize: fontSizes.md,
        fontWeight: '500',
        color: colors.black,
        marginVertical: spacing.xxl,
    },
    scrollView: {
        flex: 1,
    },
    scrollContent: {
        paddingBottom: spacing.xl,
        paddingHorizontal: spacing.lg,
    },
    hintBox: {
        alignItems: 'center',
        marginTop: 40,
    },
    hintText: {
        color: colors.gray500,
        fontSize: fontSizes.md,
    },
    emptyBox: {
        marginTop: 20,
    }
});