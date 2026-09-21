import React, { useCallback, useRef, useEffect, useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, TextInput, Keyboard, PanResponder, Image, Alert, ActivityIndicator, ScrollView } from 'react-native';
import * as ImagePicker from 'expo-image-picker';
import { uploadImageToStorage } from '../utils/image';
import {
  BottomSheetModal,
  BottomSheetView,
  BottomSheetBackdrop,
  BottomSheetBackdropProps,
  BottomSheetTextInput
} from '@gorhom/bottom-sheet';
import { Ionicons } from '@expo/vector-icons';
import { useTranslation } from 'react-i18next';
import { colors, spacing, fontSizes, borderRadius } from '../constants/theme';
import { Button } from './Button';
import { IconButton } from './IconButton';

interface ReviewFlowModalProps {
  visible: boolean;
  onClose: () => void;
  onSubmit: (rating: number, comment: string, photoUrls: string[]) => void;
}

export const ReviewFlowModal = ({ visible, onClose, onSubmit }: ReviewFlowModalProps) => {
  const { t } = useTranslation();
  const bottomSheetModalRef = useRef<BottomSheetModal>(null);
  const [rating, setRating] = useState<number>(5);
  const [showComment, setShowComment] = useState<boolean>(false);
  const [commentText, setCommentText] = useState<string>('');
  const [photos, setPhotos] = useState<{ uri: string }[]>([]);
  const [isUploading, setIsUploading] = useState<boolean>(false);
  
  const sliderWidthRef = useRef<number>(1);
  const ratingRef = useRef<number>(5);

  useEffect(() => {
    ratingRef.current = rating;
  }, [rating]);

  useEffect(() => {
    if (visible) {
      // Reset state on open
      setRating(5);
      setShowComment(false);
      setCommentText('');
      setPhotos([]);
      setIsUploading(false);
      bottomSheetModalRef.current?.present();
    } else {
      bottomSheetModalRef.current?.dismiss();
      Keyboard.dismiss();
    }
  }, [visible]);

  const handleClose = () => {
    Keyboard.dismiss();
    onClose();
  };

  const handlePickPhoto = async () => {
    if (photos.length >= 5) {
      Alert.alert(t('review.photoLimitTitle'), t('review.photoLimitDescription'));
      return;
    }

    Alert.alert(
      t('review.addPhotoTitle'),
      t('review.addPhotoDescription'),
      [
        { text: t('review.camera'), onPress: async () => {
            const { status } = await ImagePicker.requestCameraPermissionsAsync();
            if (status !== 'granted') return Alert.alert(t('review.permissionRequired'), t('review.cameraPermission'));
            const result = await ImagePicker.launchCameraAsync({
              mediaTypes: ['images'],
              allowsEditing: true,
              aspect: [4, 3],
              quality: 0.8,
            });
            if (!result.canceled && result.assets[0]) {
              setPhotos(prev => [...prev, { uri: result.assets[0].uri }]);
            }
          }
        },
        { text: t('review.gallery'), onPress: async () => {
            const { status } = await ImagePicker.requestMediaLibraryPermissionsAsync();
            if (status !== 'granted') return Alert.alert(t('review.permissionRequired'), t('review.galleryPermission'));
            const result = await ImagePicker.launchImageLibraryAsync({
              mediaTypes: ['images'],
              allowsEditing: true,
              aspect: [4, 3],
              quality: 0.8,
            });
            if (!result.canceled && result.assets[0]) {
              setPhotos(prev => [...prev, { uri: result.assets[0].uri }]);
            }
          }
        },
        { text: t('common.cancel'), style: 'cancel' }
      ]
    );
  };

  const handleSubmit = async () => {
     try {
       setIsUploading(true);
       const uploadedUrls: string[] = [];
       for (const photo of photos) {
          const url = await uploadImageToStorage(photo.uri, 'review');
          if (url) uploadedUrls.push(url);
       }
       onSubmit(rating, commentText, uploadedUrls);
       handleClose();
     } catch (e) {
       Alert.alert(t('common.error'), t('review.submitError'));
     } finally {
       setIsUploading(false);
     }
  };

  const renderBackdrop = useCallback(
    (props: BottomSheetBackdropProps) => (
      <BottomSheetBackdrop 
        {...props} 
        appearsOnIndex={0} 
        disappearsOnIndex={-1} 
        pressBehavior="close" 
      />
    ),
    []
  );

  const renderStars = () => {
    return (
      <View style={styles.starsContainer}>
        {[1, 2, 3, 4, 5].map((starIndex) => (
          <TouchableOpacity
            key={starIndex}
            activeOpacity={0.7}
            onPress={() => setRating(starIndex)}
          >
            <Ionicons
              name="star"
              size={48}
              color={starIndex <= rating ? '#FFB800' : colors.gray200}
              style={{ marginHorizontal: 4 }}
            />
          </TouchableOpacity>
        ))}
      </View>
    );
  };

  const panResponder = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => true,
      onStartShouldSetPanResponderCapture: () => true,
      onMoveShouldSetPanResponder: () => true,
      onMoveShouldSetPanResponderCapture: (evt, gestureState) => {
        return Math.abs(gestureState.dx) > 2;
      },
      onPanResponderGrant: (evt) => handlePan(evt.nativeEvent.locationX),
      onPanResponderMove: (evt) => handlePan(evt.nativeEvent.locationX),
      onPanResponderTerminationRequest: () => false,
    })
  ).current;

  const handlePan = (x: number) => {
    const ratio = Math.max(0, Math.min(1, x / sliderWidthRef.current));
    const value = Math.round(ratio * 4) + 1;
    if (value >= 1 && value <= 5) {
      setRating(value);
    }
  };

  const renderSliderIndicator = () => {
    return (
      <View style={styles.sliderWrapper}>
        <View 
           style={[styles.sliderLineContainer, { paddingVertical: 15, marginVertical: -15 }]} 
           {...panResponder.panHandlers}
           onLayout={(e) => {
             sliderWidthRef.current = e.nativeEvent.layout.width || 1;
           }}
        >
          <View style={[styles.sliderLine, styles.sliderLineBase]} pointerEvents="none" />
          {rating > 0 && (
             <View style={[
                 styles.sliderLine, 
                 styles.sliderLineActive, 
                 { width: `${((rating - 1) / 4) * 100}%` }
             ]} pointerEvents="none" />
          )}
          
          <View style={styles.sliderDotsRow} pointerEvents="none">
            {[1, 2, 3, 4, 5].map((index) => {
              const isActive = index <= rating;
              const isCurrent = index === rating;
              return (
                <View 
                    key={index} 
                    style={[
                        styles.dot,
                        isActive ? styles.dotActive : styles.dotInactive,
                        isCurrent && styles.dotCurrent
                    ]} 
                />
              );
            })}
          </View>
        </View>
        <Text style={styles.ratingLabelText}>
          {rating > 0 ? t(`review.labels.${rating}`) : ' '}
        </Text>
      </View>
    );
  };

  return (
    <BottomSheetModal
      ref={bottomSheetModalRef}
      index={0}
      onDismiss={onClose}
      backdropComponent={renderBackdrop}
      enablePanDownToClose={true}
      handleIndicatorStyle={styles.handle}
      backgroundStyle={styles.background}
      enableDynamicSizing={true}
      keyboardBehavior="interactive"
      keyboardBlurBehavior="restore"
    >
      <BottomSheetView style={styles.content}>
        <View style={styles.header}>
          <Text style={styles.title}>{t('review.title')}</Text>
          <View style={styles.closeButtonContainer}>
            <IconButton icon="close" size="small" onPress={handleClose} />
          </View>
        </View>

        <View style={styles.body}>
          <Text style={styles.headingTitle}>{t('review.headingTitle')}</Text>
          <Text style={styles.headingSubtitle}>
            {t('review.headingSubtitle')}
          </Text>

          {renderStars()}
          {renderSliderIndicator()}

          {showComment && (
            <>
              <BottomSheetTextInput
                style={styles.commentInput}
                placeholder={t('review.commentPlaceholder')}
                value={commentText}
                onChangeText={setCommentText}
                multiline
                autoFocus
                textAlignVertical="top"
              />
              <ScrollView horizontal showsHorizontalScrollIndicator={false} style={styles.photosScrollContainer}>
                {photos.map((p, idx) => (
                  <View key={idx} style={styles.photoPreviewWrapper}>
                    <Image source={{ uri: p.uri }} style={styles.photoPreview} />
                    <TouchableOpacity onPress={() => setPhotos(prev => prev.filter((_, i) => i !== idx))} style={styles.removePhotoBtn}>
                      <Ionicons name="close" size={12} color="#FFF" />
                    </TouchableOpacity>
                  </View>
                ))}
                {photos.length < 5 && (
                  <TouchableOpacity onPress={handlePickPhoto} style={styles.addPhotoBtn}>
                    <Ionicons name="camera" size={24} color={colors.gray400} />
                  </TouchableOpacity>
                )}
              </ScrollView>
            </>
          )}

          <View style={styles.footerRow}>
            {!showComment ? (
              <View style={{ flex: 1 }}>
                <Button
                  title={t('review.comment')}
                  variant="secondary"
                  onPress={() => setShowComment(true)}
                  style={styles.commentButtonOutline}
                  textStyle={styles.commentButtonText}
                />
              </View>
            ) : null}
            <View style={{ flex: 1, marginLeft: !showComment ? spacing.md : 0 }}>
              <Button
                title={isUploading ? t('review.uploading') : t('review.submit')}
                variant="primary"
                onPress={handleSubmit}
                disabled={rating === 0 || isUploading}
              />
            </View>
          </View>
        </View>
      </BottomSheetView>
    </BottomSheetModal>
  );
};

const styles = StyleSheet.create({
  background: {
    backgroundColor: colors.background,
    borderRadius: 32,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: -4 },
    shadowOpacity: 0.1,
    shadowRadius: 10,
    elevation: 20,
  },
  content: { 
    backgroundColor: colors.background, 
    borderTopLeftRadius: 32, 
    borderTopRightRadius: 32, 
    paddingBottom: spacing.xxl 
  },
  handle: { 
    width: 48, 
    height: 4, 
    backgroundColor: colors.gray400, 
    borderRadius: 2, 
    alignSelf: 'center', 
    marginTop: 8
  },
  header: { 
    flexDirection: 'row', 
    justifyContent: 'center', 
    alignItems: 'center', 
    padding: spacing.lg, 
    position: 'relative' 
  },
  title: { 
    fontSize: fontSizes.md, 
    fontWeight: '500',
    color: colors.black
  },
  closeButtonContainer: { 
    position: 'absolute', 
    right: spacing.lg, 
  },
  body: {
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    alignItems: 'center',
  },
  headingTitle: {
    fontSize: 22,
    fontWeight: '700',
    color: colors.black,
    marginBottom: spacing.xs,
    textAlign: 'center',
  },
  headingSubtitle: {
    fontSize: fontSizes.md,
    color: colors.gray600,
    textAlign: 'center',
    marginBottom: spacing.xxl,
  },
  starsContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    marginBottom: spacing.xl,
  },
  sliderWrapper: {
    width: '80%',
    alignItems: 'center',
    marginBottom: spacing.xxl,
  },
  sliderLineContainer: {
    width: '100%',
    height: 14,
    justifyContent: 'center',
  },
  sliderLine: {
    position: 'absolute',
    height: 2,
    borderRadius: 1,
    left: 0,
    right: 0, // This is for the base line
  },
  sliderLineBase: {
    backgroundColor: colors.gray200,
  },
  sliderLineActive: {
    backgroundColor: colors.primary,
    right: 'auto', // override right so width dictates it
  },
  sliderDotsRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    width: '100%',
    position: 'absolute',
  },
  dot: {
    width: 6,
    height: 6,
    borderRadius: 3,
  },
  dotInactive: {
    backgroundColor: colors.gray200,
  },
  dotActive: {
    backgroundColor: colors.primary,
  },
  dotCurrent: {
    width: 14,
    height: 14,
    borderRadius: 7,
    borderWidth: 2,
    borderColor: colors.primary,
    backgroundColor: colors.primary, // visually big solid circle
  },
  ratingLabelText: {
    marginTop: spacing.md,
    fontSize: fontSizes.md,
    fontWeight: '600',
    color: colors.primary,
  },
  commentInput: {
    width: '100%',
    backgroundColor: colors.gray100,
    borderRadius: borderRadius.md,
    padding: spacing.md,
    height: 100,
    marginBottom: spacing.xl,
    fontSize: fontSizes.md,
  },
  footerRow: {
    flexDirection: 'row',
    width: '100%',
    marginTop: spacing.md,
  },
  commentButtonOutline: {
    backgroundColor: colors.white,
    borderWidth: 1,
    borderColor: colors.primary,
  },
  commentButtonText: {
    color: colors.primary,
  },
  photosScrollContainer: {
    width: '100%',
    marginBottom: spacing.xl,
    flexDirection: 'row',
  },
  photoPreviewWrapper: {
    width: 60,
    height: 60,
    marginRight: spacing.sm,
    position: 'relative',
  },
  photoPreview: {
    width: 60,
    height: 60,
    borderRadius: borderRadius.sm,
    backgroundColor: colors.gray200,
  },
  removePhotoBtn: {
    position: 'absolute',
    top: -6,
    right: -6,
    width: 20,
    height: 20,
    borderRadius: 10,
    backgroundColor: colors.error || '#F44336',
    alignItems: 'center',
    justifyContent: 'center',
  },
  addPhotoBtn: {
    width: 60,
    height: 60,
    borderRadius: borderRadius.sm,
    borderWidth: 1,
    borderStyle: 'dashed',
    borderColor: colors.gray400,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: colors.gray100,
  },
});
