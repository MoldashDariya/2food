export interface MockRestaurant {
  id: string;
  name: string;
  description: string;
  logo?: string;
  workingHours: string;
  rating: number;
}

export const mockRestaurants: MockRestaurant[] = [
  { id: '1', name: 'Булочная №1', description: 'Свежая выпечка и кофе', workingHours: '9:00-21:00', rating: 4.9 },
  { id: '2', name: 'Sushi Wave', description: 'Роллы и поке', workingHours: '10:00-22:00', rating: 4.8 },
  { id: '3', name: 'Urban Pizza', description: 'Неаполитанская пицца', workingHours: '11:00-23:00', rating: 4.7 },
  { id: '4', name: 'Noodle House', description: 'Азиатская лапша и супы', workingHours: '9:30-21:30', rating: 4.6 },
  { id: '5', name: 'Vegan Bowl', description: 'Полезные боулы и салаты', workingHours: '8:00-20:00', rating: 4.5 },
  { id: '6', name: 'Burger Point', description: 'Бургеры и снеки', workingHours: '11:00-00:00', rating: 4.4 },
  { id: '7', name: 'Kebab Station', description: 'Гриль и шаурма', workingHours: '10:00-23:00', rating: 4.3 },
  { id: '8', name: 'Dolce Lab', description: 'Десерты и торты', workingHours: '8:30-21:00', rating: 4.9 },
  { id: '9', name: 'Taco Street', description: 'Мексиканская кухня', workingHours: '12:00-23:00', rating: 4.2 },
  { id: '10', name: 'Breakfast Club', description: 'Завтраки весь день', workingHours: '7:00-19:00', rating: 4.8 },
];
