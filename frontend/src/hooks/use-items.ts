import { createItem, CreateItemInput, deleteItems, listItems } from "@/api/items";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

export const useItems = () => {
  const queryClient = useQueryClient();

  const { data: items, refetch: queryItems } = useQuery({
    queryKey: ["items"],
    queryFn: async () => {
      const response = await listItems();
      return response.items;
    },
    enabled: true,
  });

  const handleCreateItem = useMutation({
    mutationFn: async (input: CreateItemInput) => {
      await createItem(input);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["items"] });
      toast.success("Item created successfully");
    },
    onError: () => toast.error("Error creating item"),
  });

  const handleDeleteItems = useMutation({
    mutationFn: async () => {
      await deleteItems();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["items"] });
      toast.success("Item deleted successfully");
    },
    onError: () => toast.error("Error deleting item"),
  });

  return {
    items,
    queryItems,
    handleCreateItem,
    handleDeleteItems,
  };
};
