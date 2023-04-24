import React from "react";
import {
  CommandBar,
  ICommandBarItemProps,
  ProgressIndicator,
} from "@fluentui/react";
import styles from "./CommandBarContainer.module.css";
import cn from "classnames";

const progressIndicatorStyles = {
  itemProgress: {
    padding: 0,
  },
};

const commandBarStyles = {
  root: {
    // Align the first item with the screen title.
    padding: "0 14px",
  },
};

export interface CommandBarContainerProps {
  className?: string;
  isLoading?: boolean;
  messageBar?: React.ReactNode;
  primaryItems?: ICommandBarItemProps[];
  secondaryItems?: ICommandBarItemProps[];
  children?: React.ReactNode;
  hideCommandBar?: boolean;
  headerPosition?: "static" | "sticky";
}

const CommandBarContainer: React.VFC<CommandBarContainerProps> =
  function CommandBarContainer(props) {
    const {
      className,
      isLoading,
      primaryItems,
      secondaryItems,
      messageBar,
      hideCommandBar,
      headerPosition = "sticky",
    } = props;

    return (
      <>
        <div
          className={
            headerPosition === "sticky"
              ? styles.headerSticky
              : styles.headerStatic
          }
        >
          {hideCommandBar === true ? null : (
            <CommandBar
              className={styles.commandBar}
              styles={commandBarStyles}
              items={primaryItems ?? []}
              farItems={secondaryItems}
            />
          )}
          {messageBar}
          <ProgressIndicator
            styles={progressIndicatorStyles}
            className={!isLoading ? styles.hidden : ""}
          />
        </div>
        <div className={cn(styles.content, className)}>{props.children}</div>
      </>
    );
  };

export default CommandBarContainer;
