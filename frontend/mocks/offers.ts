export interface MockOffer {
  id: string;
  restaurantId: string;
  title: string;
  description: string;
  image?: string;
  originalPrice: number;
  discountedPrice: number;
  readyInMinutes: number;
  supplyLeft: number;
}

export const mockOffers: MockOffer[] = [
  {
    id: 'offer-1',
    restaurantId: '1',
    title: 'Булочка',
    description: 'Круассан + улитка + мини-кофе',
    image: 'https://images.unsplash.com/photo-1509440159596-0249088772ff?auto=format&fit=crop&w=900&q=80',
    originalPrice: 2800,
    discountedPrice: 1500,
    readyInMinutes: 25,
    supplyLeft: 8,
  },
  {
    id: 'offer-2',
    restaurantId: '2',
    title: 'Sushi Surprise Box',
    description: 'Микс роллов и закусок',
    image: 'https://images.unsplash.com/photo-1579871494447-9811cf80d66c?auto=format&fit=crop&w=900&q=80',
    originalPrice: 5200,
    discountedPrice: 2600,
    readyInMinutes: 40,
    supplyLeft: 5,
  },
  {
    id: 'offer-3',
    restaurantId: '3',
    title: 'Pizza Combo',
    description: '2 кусочка пиццы + напиток',
    image: 'https://images.unsplash.com/photo-1513104890138-7c749659a591?auto=format&fit=crop&w=900&q=80',
    originalPrice: 3500,
    discountedPrice: 1900,
    readyInMinutes: 30,
    supplyLeft: 12,
  },
  {
    id: 'offer-4',
    restaurantId: '4',
    title: 'Noodle Bowl Saver',
    description: 'Лапша + спринг-ролл',
    image: 'https://images.unsplash.com/photo-1612929633738-8fe44f7ec841?auto=format&fit=crop&w=900&q=80',
    originalPrice: 3400,
    discountedPrice: 1700,
    readyInMinutes: 20,
    supplyLeft: 6,
  },
  {
    id: 'offer-5',
    restaurantId: '5',
    title: 'Vegan Lunch Pack',
    description: 'Боул + суп + лимонад',
    image: 'https://images.unsplash.com/photo-1547592166-23ac45744acd?auto=format&fit=crop&w=900&q=80',
    originalPrice: 4200,
    discountedPrice: 2200,
    readyInMinutes: 25,
    supplyLeft: 7,
  },
  {
    id: 'offer-6',
    restaurantId: '6',
    title: 'Burger Deal',
    description: 'Бургер + картофель',
    image: 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?auto=format&fit=crop&w=900&q=80',
    originalPrice: 3900,
    discountedPrice: 2100,
    readyInMinutes: 15,
    supplyLeft: 9,
  },
  {
    id: 'offer-7',
    restaurantId: '7',
    title: 'Grill Box',
    description: 'Мясной сет с гарниром',
    image: 'https://images.unsplash.com/photo-1529193591184-b1d58069ecdd?auto=format&fit=crop&w=900&q=80',
    originalPrice: 4800,
    discountedPrice: 2500,
    readyInMinutes: 35,
    supplyLeft: 4,
  },
  {
    id: 'offer-8',
    restaurantId: '8',
    title: 'Dessert Duo',
    description: '2 десерта по выбору',
    image: 'https://images.unsplash.com/photo-1488477181946-6428a0291777?auto=format&fit=crop&w=900&q=80',
    originalPrice: 3000,
    discountedPrice: 1500,
    readyInMinutes: 20,
    supplyLeft: 10,
  },
  {
    id: 'offer-9',
    restaurantId: '9',
    title: 'Taco Night',
    description: '3 тако + соус',
    image: 'https://images.unsplash.com/photo-1613514785940-daed07799d9b?auto=format&fit=crop&w=900&q=80',
    originalPrice: 3600,
    discountedPrice: 1800,
    readyInMinutes: 18,
    supplyLeft: 11,
  },
  {
    id: 'offer-10',
    restaurantId: '10',
    title: 'Breakfast Bag',
    description: 'Сэндвич + кофе + фрукт',
    image: 'https://images.unsplash.com/photo-1533089860892-a7c6f0a88666?auto=format&fit=crop&w=900&q=80',
    originalPrice: 3200,
    discountedPrice: 1600,
    readyInMinutes: 10,
    supplyLeft: 13,
  },
];
