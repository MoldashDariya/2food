import { mockOffers } from './offers';
import { mockRestaurants, type MockRestaurant } from './restaurants';

interface SearchRestaurantSource {
  id: string;
  name: string;
  description?: string;
  logo?: string;
  image_url?: string;
  workingHours?: string;
  rating?: number;
}

export interface SearchRestaurantEntity {
  id: string;
  name: string;
  description: string;
  logo?: string;
  workingHours: string;
  rating: number;
  offersCount: number;
  totalSupplyLeft: number;
  allProducts: any[]; 
  filteredProducts: any[];
  topOfferTitle?: string;
  bestOffer?: {
    id: string;
    title: string;
    image?: string;
    originalPrice: number;
    discountedPrice: number;
    readyInMinutes: number;
    supplyLeft: number;
  };
}

const toSearchEntity = (
  restaurant: SearchRestaurantSource,
  index: number,
): SearchRestaurantEntity => {
  const matchedOffers = mockOffers.filter((offer) => offer.restaurantId === restaurant.id);
  const offers =
    matchedOffers.length > 0
      ? matchedOffers
      : [mockOffers[index % mockOffers.length]];
  const bestOffer =
    offers.length > 0
      ? [...offers].sort((a, b) => a.discountedPrice - b.discountedPrice)[0]
      : undefined;
  const totalSupplyLeft = offers.reduce((sum, offer) => sum + offer.supplyLeft, 0);
  const fallback = mockRestaurants.find((item) => item.id === restaurant.id);

  return {
    id: restaurant.id,
    name: restaurant.name,
    description: restaurant.description ?? 'Без описания',
    logo: restaurant.logo ?? restaurant.image_url ?? fallback?.logo,
    workingHours: restaurant.workingHours ?? fallback?.workingHours ?? '9:00-21:00',
    rating: restaurant.rating ?? fallback?.rating ?? 4.5,
    offersCount: offers.length,
    totalSupplyLeft,
    allProducts: offers, 
    filteredProducts: offers,
    topOfferTitle: bestOffer?.title,
    bestOffer: bestOffer
      ? {
          id: bestOffer.id,
          title: bestOffer.title,
          image: bestOffer.image,
          originalPrice: bestOffer.originalPrice,
          discountedPrice: bestOffer.discountedPrice,
          readyInMinutes: bestOffer.readyInMinutes,
          supplyLeft: bestOffer.supplyLeft,
        }
      : undefined,
  };
};

export const getSearchRestaurants = (
  restaurants?: SearchRestaurantSource[],
): SearchRestaurantEntity[] => {
  const source: SearchRestaurantSource[] = restaurants ?? mockRestaurants;
  return source.map((restaurant, index) => toSearchEntity(restaurant, index));
};

export const searchRestaurants = (
  query: string,
  restaurants?: SearchRestaurantSource[],
): SearchRestaurantEntity[] => {
  const normalizedQuery = query.trim().toLowerCase();
  const entities = getSearchRestaurants(restaurants);

  if (!normalizedQuery) {
    return entities;
  }

  return entities.filter((entity) => {
    const name = entity.name.toLowerCase();
    const description = entity.description.toLowerCase();
    const offerTitle = entity.topOfferTitle?.toLowerCase() ?? '';
    return (
      name.includes(normalizedQuery) ||
      description.includes(normalizedQuery) ||
      offerTitle.includes(normalizedQuery)
    );
  });
};

export type SortOption = 'default' | 'rating' | 'discount' | 'price' | 'time';
export type PriceOption = 'all' | 'under2000' | 'under3000';
export type TimeOption = 'all' | 'under20' | 'under30';

interface FilterOptions {
  sort: SortOption;
  price: PriceOption;
  time: TimeOption;
}

export const applySearchFilters = (
  entities: SearchRestaurantEntity[],
  options: FilterOptions,
  customPriceRange?: [number, number],
): SearchRestaurantEntity[] => {
  const filtered = entities.filter((entity) => {
    const currentPrice = entity.bestOffer?.discountedPrice ?? Number.MAX_SAFE_INTEGER;
    const currentTime = entity.bestOffer?.readyInMinutes ?? Number.MAX_SAFE_INTEGER;

    let priceOk = true;
    
    if (customPriceRange) {
      priceOk = currentPrice >= customPriceRange[0] && currentPrice <= customPriceRange[1];
    } else {
      priceOk =
        options.price === 'all' ||
        (options.price === 'under2000' && currentPrice <= 2000) ||
        (options.price === 'under3000' && currentPrice <= 3000);
    }

    const timeOk =
      options.time === 'all' ||
      (options.time === 'under20' && currentTime <= 20) ||
      (options.time === 'under30' && currentTime <= 30);

    return priceOk && timeOk;
  });

  const sorted = [...filtered];

  if (options.sort === 'rating') {
    sorted.sort((a, b) => b.rating - a.rating);
  } 
  else if (options.sort === 'discount') {
    sorted.sort((a, b) => {
      const left = (a.bestOffer?.originalPrice ?? 0) - (a.bestOffer?.discountedPrice ?? 0);
      const right = (b.bestOffer?.originalPrice ?? 0) - (b.bestOffer?.discountedPrice ?? 0);
      return right - left;
    });
  }
  else if (options.sort === 'price') {
    sorted.sort((a, b) => 
      (a.bestOffer?.discountedPrice ?? 0) - (b.bestOffer?.discountedPrice ?? 0)
    );
  }
  else if (options.sort === 'time') {
    sorted.sort((a, b) => 
      (a.bestOffer?.readyInMinutes ?? 0) - (b.bestOffer?.readyInMinutes ?? 0)
    );
  }

  return sorted;
};