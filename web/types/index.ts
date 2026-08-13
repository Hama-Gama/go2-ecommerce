export interface Product {
  id: number;
  category_id: number;
  title: string;
  slug: string;
  description: string;
  price: number;
  stock: number;
  created_at: string;
  updated_at: string;
}

export interface CartItem {
  product: Product;
  quantity: number;
}

export interface CreateOrderRequest {
  user_id: number;
  items: {
    product_id: number;
    quantity: number;
  }[];
}

export interface OrderResponse {
  id: number;
  user_id: number;
  status: string;
  total_amount: number;
}