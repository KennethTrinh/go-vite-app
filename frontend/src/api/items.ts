import { axiosClient } from "@/lib/axios-client";
import axios from "axios";

// curl -X POST http://localhost:8000/items -H "Content-Type: application/json" -d '{"name":"test","description":"test","icon":"❤️","color":"test","time":1}'
// curl http://localhost:8000/items
// curl -X DELETE http://localhost:8000/items/1

export type CreateItemInput = {
  name: string;
  description: string;
  icon: string;
  color: string;
  time: number;
};

export const createItem = async (input: CreateItemInput) => {
  const url = "/items";
  await axiosClient.post(url, input);
};

export const listItems = async () => {
  const url = "/items";
  try {
    const response = await axiosClient.get(url);
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error) && error.response?.status === 429) {
      const retryAfter = error.response.headers["retry-after"];
      return {
        items: [],
        error: {
          message: `Rate limited, retry after ${retryAfter} seconds`,
        },
      };
    }
    return {
      items: [],
      error: {
        message: (error as any).response?.data?.errors[0].message,
      },
    };
  }
};

export const deleteItems = async () => {
  const url = `/items`;
  await axiosClient.delete(url);
};
