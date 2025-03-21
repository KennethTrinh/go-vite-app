import { cn } from "@/lib/utils";
import { AnimatedList, Item, Notification } from "./animated-list";
import { Link } from "react-router";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useItems } from "@/hooks/use-items";

export const Component = () => {
  const {
    items,
    handleCreateItem: { mutate: createItem },
    handleDeleteItems: { mutate: deleteItems },
  } = useItems();

  return (
    <div
      className={cn(
        "relative flex min-h-screen w-full flex-col p-6 overflow-hidden rounded-lg border bg-background "
      )}
    >
      <Link
        to="/"
        className="flex items-center text-sm font-medium text-muted-foreground hover:text-gray-900 self-start"
      >
        <ArrowLeft className="mr-2 h-5 w-5" />
        Back to Homepage
      </Link>
      <div className="flex items-center justify-between min-w-[400px] mx-auto my-6 gap-x-2">
        <Button
          className="flex items-center justify-center w-1/2 h-12 "
          onClick={() => {
            createItem({
              name: "Payment received",
              description: "Magic UI",
              time: 15,
              icon: "💸",
              color: "#00C9A7",
            });
            createItem({
              name: "User signed up",
              description: "Magic UI",
              time: 10,
              icon: "👤",
              color: "#FFB800",
            });
            createItem({
              name: "New message",
              description: "Magic UI",
              time: 5,
              icon: "💬",
              color: "#FF3D71",
            });
            createItem({
              name: "New event",
              description: "Magic UI",
              time: 2,
              icon: "🗞️",
              color: "#1E86FF",
            });
          }}
        >
          Create 4 items
        </Button>
        <Button
          className="flex items-center justify-center w-1/2 h-12 "
          variant="destructive"
          onClick={() => {
            deleteItems();
          }}
        >
          Delete Items
        </Button>
      </div>
      <AnimatedList>
        {items?.map((item: Item, idx: number) => (
          <Notification {...item} key={idx} />
        ))}
      </AnimatedList>
    </div>
  );
};
