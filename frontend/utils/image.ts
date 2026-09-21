import * as ImagePicker from 'expo-image-picker';
import { getPresignedUrl } from '../services/admin';

export const clamp = (v: number, min: number, max: number) => Math.min(max, Math.max(min, v));

export const getAdaptiveQuality = (
  width?: number,
  height?: number,
  currentBytes?: number,
): number => {
  const TARGET = 800 * 1024;
  const MIN_Q = 0.05;
  const MAX_Q = 0.9;

  if (currentBytes && currentBytes <= TARGET) return MAX_Q;

  if (currentBytes && currentBytes > 0) {
    const ratio = Math.pow(TARGET / currentBytes, 0.85);
    return Number(clamp(ratio, MIN_Q, MAX_Q).toFixed(2));
  }

  if (!width || !height) return 0.25;
  const pixels = width * height;
  const estimated = TARGET / (pixels * 1.25);
  return Number(clamp(estimated, MIN_Q, MAX_Q).toFixed(2));
};

export const uploadImageToStorage = async (localUri: string, target: 'avatar' | 'restaurant_logo' | 'restaurant_banner' | 'product' | 'category' | 'review') => {
  try {
    const response = await fetch(localUri);
    const blob = await response.blob();
    const filename = `${target}_${Date.now()}.jpg`;

    const { upload_url, file_url } = await getPresignedUrl({
      content_type: 'image/jpeg',
      target,
    });

    const uploadResponse = await fetch(upload_url, {
      method: 'PUT',
      body: blob,
      headers: { 'Content-Type': 'image/jpeg' },
    });

    if (!uploadResponse.ok) return null;
    return file_url;
  } catch (e) {
    console.error(`Error uploading ${target} image:`, e);
    return null;
  }
};
