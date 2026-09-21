import React, { useState, useRef, useMemo } from 'react';
import {
    StyleSheet,
    View,
    Text,
    FlatList,
    Image,
    Dimensions,
    NativeSyntheticEvent,
    NativeScrollEvent,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { Button } from '../components';
import { colors, fontSizes, spacing } from '../constants/theme';
import { Ionicons } from '@expo/vector-icons';
import { useAuthStore } from '../stores/authStore';
import onboarding_01 from '../assets/pictures/onboarding_01.png';
import onboarding_03 from '../assets/pictures/onboarding_03.png';

const { width, height } = Dimensions.get('window');

export default function OnboardingScreen() {
    const { t } = useTranslation();
    const [currentSlideIndex, setCurrentSlideIndex] = useState(0);
    const flatListRef = useRef<FlatList>(null);
    const router = useRouter();
    const setHasOnboarded = useAuthStore((s) => s.setHasOnboarded);

    const SLIDES = useMemo(() => ([
        {
            id: '1',
            title: t('onboarding.slide1Title'),
            description: t('onboarding.slide1Description'),
            image: onboarding_01,
        },
        {
            id: '2',
            title: t('onboarding.slide2Title'),
            subtitle: t('onboarding.slide2Subtitle'),
            steps: [
                {
                    icon: 'phone-portrait-outline',
                    text: t('onboarding.step1Title'),
                    subtext: t('onboarding.step1Description'),
                },
                {
                    icon: 'wallet-outline',
                    text: t('onboarding.step2Title'),
                    subtext: t('onboarding.step2Description'),
                },
                {
                    icon: 'storefront-outline',
                    text: t('onboarding.step3Title'),
                    subtext: t('onboarding.step3Description'),
                },
            ],
        },
        {
            id: '3',
            title: t('onboarding.slide3Title'),
            description: t('onboarding.slide3Description'),
            image: onboarding_03,
        },
    ]), [t]);

    const updateCurrentSlideIndex = (e: NativeSyntheticEvent<NativeScrollEvent>) => {
        const contentOffsetX = e.nativeEvent.contentOffset.x;
        const currentIndex = Math.round(contentOffsetX / width);
        setCurrentSlideIndex(currentIndex);
    };

    const handleNext = () => {
        const nextSlideIndex = currentSlideIndex + 1;
        if (nextSlideIndex < SLIDES.length) {
            flatListRef.current?.scrollToOffset({
                offset: nextSlideIndex * width,
                animated: true,
            });
            setCurrentSlideIndex(nextSlideIndex);
        } else {
            setHasOnboarded().then(() => {
                router.replace('/login');
            });
        }
    };

    const renderSteps = (steps: any[]) => (
        <View style={styles.stepsContainer}>
            {steps.map((step, index) => (
                <View key={index} style={styles.stepRow}>
                    {index < steps.length - 1 && <View style={styles.stepLine} />}
                    <View style={styles.iconCircle}>
                        <Ionicons name={step.icon} size={24} color={colors.white} />
                    </View>
                    <View style={styles.stepTextContent}>
                        <Text style={styles.stepTitle}>{step.text}</Text>
                        <Text style={styles.stepSubtext}>{step.subtext}</Text>
                    </View>
                </View>
            ))}
        </View>
    );

    const renderSlide = ({ item }: { item: any }) => {
        if (item.id === '2') {
            return (
                <View style={styles.slide}>
                    <View style={styles.stepsContentWrapper}>
                        <Text style={styles.title}>{item.title}</Text>
                        <Text style={styles.slideSubtitle}>{item.subtitle}</Text>
                        {renderSteps(item.steps)}
                    </View>
                </View>
            );
        }

        return (
            <View style={styles.slide}>
                <View style={styles.imageContainer}>
                    <Image
                        source={item.image}
                        style={styles.image}
                        resizeMode="contain"
                    />
                </View>

                <View style={styles.textContainer}>
                    <Text style={styles.title}>{item.title}</Text>
                    {item.description && <Text style={styles.description}>{item.description}</Text>}
                </View>
            </View>
        );
    };

    return (
        <SafeAreaView style={styles.container}>
            <FlatList
                ref={flatListRef}
                data={SLIDES}
                renderItem={renderSlide}
                horizontal
                pagingEnabled
                showsHorizontalScrollIndicator={false}
                onMomentumScrollEnd={updateCurrentSlideIndex}
                keyExtractor={(item) => item.id}
                bounces={false}
            />

            <View style={styles.footer}>
                <View style={styles.indicatorContainer}>
                    {SLIDES.map((_, index) => (
                        <View
                            key={index}
                            style={[
                                styles.indicator,
                                currentSlideIndex === index && styles.activeIndicator,
                            ]}
                        />
                    ))}
                </View>

                <View style={styles.buttonWrapper}>
                    <Button
                        title={t('onboarding.continue')}
                        onPress={handleNext}
                    />
                </View>
            </View>
        </SafeAreaView>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: colors.white,
    },
    slide: {
        width: width,
        height: '100%',
        alignItems: 'center',
        justifyContent: 'center',
        paddingHorizontal: spacing.lg,
    },
    imageContainer: {
        alignItems: 'center',
        justifyContent: 'center',
        marginBottom: 32,
    },
    image: {
        width: 300,
        height: 200,
    },
    textContainer: {
        width: '100%',
        alignItems: 'center',
    },
    title: {
        fontSize: fontSizes.xxxl,
        fontWeight: '600',
        color: colors.black,
        textAlign: 'center',
        marginBottom: spacing.md,
    },
    description: {
        fontSize: fontSizes.md,
        color: colors.black,
        textAlign: 'center',
        lineHeight: 20,
    },
    stepsContentWrapper: {
        width: '100%',
    },
    slideSubtitle: {
        fontSize: 16,
        textAlign: 'center',
        marginBottom: 40,
        color: colors.black,
        lineHeight: 22,
    },
    stepsContainer: {
        width: '100%',
    },
    stepRow: {
        flexDirection: 'row',
        marginBottom: 56,
        alignItems: 'center',
        position: 'relative',
    },
    iconCircle: {
        width: 50,
        height: 50,
        borderRadius: 25,
        backgroundColor: colors.primary || '#4CAF50',
        justifyContent: 'center',
        alignItems: 'center',
        zIndex: 2,
        shadowColor: "#b4ebb5",
        shadowOffset: {
            width: 0,
            height: 3,
        },
        shadowOpacity: 5,
        shadowRadius: 4,
        elevation: 5,
    },
    stepTextContent: {
        marginLeft: spacing.xxl,
        flex: 1,
    },
    stepTitle: {
        fontSize: fontSizes.md,
        fontWeight: '600',
        color: colors.black,
    },
    stepSubtext: {
        fontSize: fontSizes.xs,
        color: colors.black,
        marginTop: spacing.xs,
        lineHeight: 16,
    },
    stepLine: {
        position: 'absolute',
        left: 24,
        top: 50,
        bottom: -65,
        width: 1.5,
        backgroundColor: colors.gray400,
        zIndex: 1,
    },
    footer: {
        paddingHorizontal: spacing.lg,

    },
    indicatorContainer: {
        flexDirection: 'row',
        justifyContent: 'center',
        marginBottom: spacing.xl,
    },
    indicator: {
        height: 12,
        width: 12,
        borderRadius: 8,
        backgroundColor: '#c7f0c7',
        marginHorizontal: 5,
    },
    activeIndicator: {
        backgroundColor: colors.primary,
    },
    buttonWrapper: {
        width: '100%',
    },
});
