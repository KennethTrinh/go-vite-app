import { cn } from "@/lib/utils";
import React, { ReactElement, useEffect, useMemo, useState, useRef } from "react";

export interface AnimatedListProps {
  className?: string;
  children: React.ReactNode;
  delay?: number;
}
export interface Item {
  name: string;
  description: string;
  icon: string;
  color: string;
  time: string;
}

export const AnimatedList = React.memo(
  ({ className, children, delay = 200 }: AnimatedListProps) => {
    const [index, setIndex] = useState(0);
    const childrenArray = React.Children.toArray(children);
    const intervalRef = useRef<NodeJS.Timeout | null>(null);

    useEffect(() => {
      // Only create interval if we haven't shown all items yet
      if (index < childrenArray.length) {
        intervalRef.current = setInterval(() => {
          setIndex((prevIndex) => {
            const nextIndex = prevIndex + 1;

            // If we've reached the end, clear the interval
            if (nextIndex >= childrenArray.length) {
              if (intervalRef.current) {
                clearInterval(intervalRef.current);
              }
              return childrenArray.length; // Cap at final length
            }

            return nextIndex;
          });
        }, delay);
      }

      // Clean up interval when component unmounts or dependencies change
      return () => {
        if (intervalRef.current) {
          clearInterval(intervalRef.current);
          intervalRef.current = null;
        }
      };
    }, [childrenArray.length, delay, index]);

    const itemsToShow = useMemo(
      () => childrenArray.slice(0, Math.min(index + 1, childrenArray.length)).reverse(),
      [index, childrenArray]
    );

    return (
      <div className={`flex flex-col items-center gap-4 ${className}`}>
        {itemsToShow.map((item) => (
          <AnimatedListItem key={(item as ReactElement).key}>{item}</AnimatedListItem>
        ))}
      </div>
    );
  }
);

AnimatedList.displayName = "AnimatedList";

export function AnimatedListItem({ children }: { children: React.ReactNode }) {
  const [isVisible, setIsVisible] = useState(false);
  const itemRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Add animation after mount to trigger the transition
    const timer = setTimeout(() => {
      setIsVisible(true);
    }, 10);

    return () => clearTimeout(timer);
  }, []);

  return (
    <div
      ref={itemRef}
      className={`mx-auto w-full transition-all duration-300 origin-top ${
        isVisible ? "opacity-100 scale-100" : "opacity-0 scale-0"
      }`}
      style={{
        transitionTimingFunction: "cubic-bezier(0.2, 0.8, 0.2, 1.2)",
      }}
    >
      {children}
    </div>
  );
}

export const Notification = ({ name, description, icon, color, time }: Item) => {
  return (
    <figure
      className={cn(
        "relative mx-auto min-h-fit w-full max-w-[400px] cursor-pointer overflow-hidden rounded-2xl p-4",
        // animation styles
        "transition-all duration-200 ease-in-out hover:scale-[103%]",
        // light styles
        "bg-white [box-shadow:0_0_0_1px_rgba(0,0,0,.03),0_2px_4px_rgba(0,0,0,.05),0_12px_24px_rgba(0,0,0,.05)]",
        // dark styles
        "transform-gpu dark:bg-transparent dark:backdrop-blur-md dark:[border:1px_solid_rgba(255,255,255,.1)] dark:[box-shadow:0_-20px_80px_-20px_#ffffff1f_inset]"
      )}
    >
      <div className="flex flex-row items-center gap-3">
        <div
          className="flex size-10 items-center justify-center rounded-2xl"
          style={{
            backgroundColor: color,
          }}
        >
          <span className="text-lg">{icon}</span>
        </div>
        <div className="flex flex-col overflow-hidden">
          <figcaption className="flex flex-row items-center whitespace-pre text-lg font-medium dark:text-white ">
            <span className="text-sm sm:text-lg">{name}</span>
            <span className="mx-1">·</span>
            <span className="text-xs text-gray-500">{time}m ago</span>
          </figcaption>
          <p className="text-sm font-normal dark:text-white/60">{description}</p>
        </div>
      </div>
    </figure>
  );
};
