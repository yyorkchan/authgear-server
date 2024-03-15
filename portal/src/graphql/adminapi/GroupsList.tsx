import React, { useCallback, useContext, useMemo, useState } from "react";
import cn from "classnames";
import { GroupsListFragment } from "./query/groupsListQuery.generated";
import useDelayedValue from "../../hook/useDelayedValue";
import {
  ColumnActionsMode,
  DetailsRow,
  IColumn,
  IDetailsRowProps,
} from "@fluentui/react";
import styles from "./GroupsList.module.css";
import { useSystemConfig } from "../../context/SystemConfigContext";
import { useParams } from "react-router-dom";
import { Context } from "@oursky/react-messageformat";
import Link from "../../Link";
import DeleteGroupDialog, {
  DeleteGroupDialogData,
} from "../../components/roles-and-groups/DeleteGroupDialog";
import DescriptionCell from "../../components/roles-and-groups/list/DescriptionCell";
import ActionButtonCell from "../../components/roles-and-groups/list/ActionButtonCell";
import TextCell from "../../components/roles-and-groups/list/TextCell";
import RolesAndGroupsBaseList from "../../components/roles-and-groups/list/RolesAndGroupsBaseList";

interface GroupsListProps {
  className?: string;
  isSearch: boolean;
  loading: boolean;
  groups: GroupsListFragment | null;
  offset: number;
  pageSize: number;
  totalCount?: number;
  onChangeOffset?: (offset: number) => void;
}

interface GroupListItem {
  id: string;
  key: string;
  name: string | null;
  description: string | null;
}

const isGroupListItem = (value: unknown): value is GroupListItem => {
  if (!(value instanceof Object)) {
    return false;
  }
  return "key" in value && "id" in value;
};

const GroupsList: React.VFC<GroupsListProps> = function GroupsList(props) {
  const {
    className,
    loading: rawLoading,
    isSearch,
    offset,
    pageSize,
    totalCount,
    onChangeOffset,
  } = props;
  const edges = props.groups?.edges;
  const loading = useDelayedValue(rawLoading, 500);
  const { renderToString } = useContext(Context);
  const { themes } = useSystemConfig();
  const { appID } = useParams() as { appID: string };
  const columns: IColumn[] = [
    {
      key: "name",
      fieldName: "name",
      name: renderToString("GroupsList.column.name"),
      minWidth: 150,
      maxWidth: 260,
      columnActionsMode: ColumnActionsMode.disabled,
    },
    {
      key: "key",
      fieldName: "key",
      name: renderToString("GroupsList.column.key"),
      minWidth: 150,
      maxWidth: 260,
      columnActionsMode: ColumnActionsMode.disabled,
    },
    {
      key: "description",
      fieldName: "description",
      name: renderToString("GroupsList.column.description"),
      minWidth: 300,
      columnActionsMode: ColumnActionsMode.disabled,
    },
    {
      key: "action",
      fieldName: "action",
      name: renderToString("GroupsList.column.action"),
      minWidth: 77,
      maxWidth: 77,
      columnActionsMode: ColumnActionsMode.disabled,
    },
  ];
  const items: GroupListItem[] = useMemo(() => {
    const items = [];
    if (edges != null) {
      for (const edge of edges) {
        const node = edge?.node;
        if (node != null) {
          items.push({
            id: node.id,
            name: node.name ?? null,
            key: node.key,
            description: node.description ?? null,
          });
        }
      }
    }
    return items;
  }, [edges]);

  const onRenderGroupRow = React.useCallback(
    (props?: IDetailsRowProps) => {
      if (props == null) {
        return null;
      }
      const targetPath = isGroupListItem(props.item)
        ? `/project/${appID}/user-management/groups/${props.item.id}/details`
        : ".";
      return (
        <Link to={targetPath} className="contents">
          <DetailsRow {...props} />
        </Link>
      );
    },
    [appID]
  );
  const [deleteGroupDialogData, setDeleteGroupDialogData] =
    useState<DeleteGroupDialogData | null>(null);
  const onClickDeleteGroup = useCallback(
    (e: React.MouseEvent<unknown>, item: GroupListItem) => {
      e.preventDefault();
      e.stopPropagation();
      setDeleteGroupDialogData({
        groupID: item.id,
        groupName: item.name,
        groupKey: item.key,
      });
    },
    []
  );
  const dismissDeleteGroupDialog = useCallback(() => {
    setDeleteGroupDialogData(null);
  }, []);

  const onRenderGroupItemColumn = useCallback(
    (item: GroupListItem, _index?: number, column?: IColumn) => {
      switch (column?.key) {
        case "description":
          return (
            <DescriptionCell>
              {item[column.key as keyof GroupListItem] ?? ""}
            </DescriptionCell>
          );
        case "action": {
          return (
            <ActionButtonCell
              text={renderToString("GroupsList.delete-group")}
              onClick={(e) => {
                onClickDeleteGroup(e, item);
              }}
            />
          );
        }
        default:
          return (
            <TextCell>
              {item[column?.key as keyof GroupListItem] ?? ""}
            </TextCell>
          );
      }
    },
    [themes.destructive, onClickDeleteGroup]
  );

  const paginationProps = useMemo(
    () => ({
      isSearch,
      offset,
      pageSize,
      totalCount,
      onChangeOffset,
    }),
    []
  );

  return (
    <>
      <div className={cn(styles.root, className)}>
        <RolesAndGroupsBaseList
          loading={loading}
          onRenderRow={onRenderGroupRow}
          onRenderItemColumn={onRenderGroupItemColumn}
          items={items}
          columns={columns}
          pagination={paginationProps}
        />
        <DeleteGroupDialog
          onDismiss={dismissDeleteGroupDialog}
          data={deleteGroupDialogData}
        />
      </div>
    </>
  );
};

export default GroupsList;
