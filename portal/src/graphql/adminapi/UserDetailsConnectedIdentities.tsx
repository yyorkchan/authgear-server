import React, {
  useMemo,
  useCallback,
  useContext,
  useState,
  useEffect,
} from "react";
import cn from "classnames";
import { useNavigate } from "react-router-dom";
import { FormattedMessage, Context } from "@oursky/react-messageformat";
import {
  DefaultButton,
  Dialog,
  DialogFooter,
  Icon,
  IContextualMenuProps,
  List,
  PrimaryButton,
  Text,
} from "@fluentui/react";

import PrimaryIdentitiesSelectionForm from "./PrimaryIdentitiesSelectionForm";
import ButtonWithLoading from "../../ButtonWithLoading";
import ListCellLayout from "../../ListCellLayout";
import { useDeleteIdentityMutation } from "./mutations/deleteIdentityMutation";
import { formatDatetime } from "../../util/formatDatetime";
import { parseError } from "../../util/error";
import { Violation } from "../../util/validation";
import { OAuthSSOProviderType } from "../../types";
import { destructiveTheme, verifyButtonTheme } from "../../theme";

import styles from "./UserDetailsConnectedIdentities.module.scss";

interface IdentityClaim extends Record<string, unknown> {
  email?: string;
  phone_number?: string;
  preferred_username?: string;
  "https://authgear.com/claims/oauth/provider_type"?: OAuthSSOProviderType;
  "https://authgear.com/claims/login_id/type"?: LoginIDIdentityType;
}

interface Identity {
  id: string;
  type: "ANONYMOUS" | "LOGIN_ID" | "OAUTH";
  claims: IdentityClaim;
  createdAt: string;
  updatedAt: string;
}

interface UserDetailsConnectedIdentitiesProps {
  identities: Identity[];
  availableLoginIdIdentities: string[];
}

const loginIdIdentityTypes = ["email", "phone", "username"] as const;
type LoginIDIdentityType = typeof loginIdIdentityTypes[number];
type IdentityType = LoginIDIdentityType | "oauth";

interface OAuthIdentityListItem {
  id: string;
  providerType?: OAuthSSOProviderType;
  name?: string;
  verified: boolean;
  connectedOn: string;
}

interface EmailIdentityListItem {
  id: string;
  email?: string;
  verified: boolean;
  addedOn: string;
}

interface PhoneIdentityListItem {
  id: string;
  phone?: string;
  verified: boolean;
  addedOn: string;
}

interface UsernameIdentityListItem {
  id: string;
  username?: string;
  addedOn: string;
}

export interface IdentityLists {
  oauth: OAuthIdentityListItem[];
  email: EmailIdentityListItem[];
  phone: PhoneIdentityListItem[];
  username: UsernameIdentityListItem[];
}

interface IdentityListCellProps {
  identityID: string;
  oauthProviderType?: OAuthSSOProviderType;
  identityType: IdentityType;
  identityName?: string;
  addedOn?: string;
  connectedOn?: string;
  verified?: boolean;
  toggleVerified?: (identityID: string, verified: boolean) => void;
  onRemoveClicked: (identityID: string, identityName: string) => void;
}

interface VerifyButtonProps {
  verified: boolean;
  toggleVerified: (verified: boolean) => void;
}

interface ConfirmationDialogData {
  identityID: string;
  identityName: string;
}

interface ErrorDialogData {
  message: string;
}

const oauthIconMap: Record<OAuthSSOProviderType, React.ReactNode> = {
  apple: <i className={cn("fab", "fa-apple", styles.widgetLabelIcon)} />,
  google: <i className={cn("fab", "fa-google", styles.widgetLabelIcon)} />,
  facebook: <i className={cn("fab", "fa-facebook", styles.widgetLabelIcon)} />,
  linkedin: <i className={cn("fab", "fa-linkedin", styles.widgetLabelIcon)} />,
  azureadv2: (
    <i className={cn("fab", "fa-microsoft", styles.widgetLabelIcon)} />
  ),
};

const loginIdIconMap: Record<LoginIDIdentityType, React.ReactNode> = {
  email: <Icon iconName="Mail" />,
  phone: <Icon iconName="CellPhone" />,
  username: <Icon iconName="Accounts" />,
};

const removeButtonTextId: Record<IdentityType, "remove" | "disconnect"> = {
  oauth: "disconnect",
  email: "remove",
  phone: "remove",
  username: "remove",
};

function getIcon(
  identityType: IdentityType,
  providerType?: OAuthSSOProviderType
) {
  if (identityType === "oauth") {
    if (providerType != null) {
      return oauthIconMap[providerType];
    }
    return null;
  }
  return loginIdIconMap[identityType];
}

function getErrorMessageIdsFromViolation(violations: Violation[]) {
  const errorMessageIds: string[] = [];
  const unknownViolations: Violation[] = [];
  for (const violation of violations) {
    switch (violation.kind) {
      case "RemoveLastIdentity":
        errorMessageIds.push(
          "UserDetails.connected-identities.remove-identity-error.connot-remove-last"
        );
        break;
      default:
        unknownViolations.push(violation);
        break;
    }
  }
  return { errorMessageIds, unknownViolations };
}

const VerifyButton: React.FC<VerifyButtonProps> = function VerifyButton(
  props: VerifyButtonProps
) {
  const { verified, toggleVerified } = props;

  const onClickVerify = useCallback(() => {
    toggleVerified(true);
  }, [toggleVerified]);

  const onClickUnverify = useCallback(() => {
    toggleVerified(false);
  }, [toggleVerified]);

  if (verified) {
    return (
      <DefaultButton
        className={cn(styles.controlButton, styles.unverifyButton)}
        onClick={onClickUnverify}
      >
        <FormattedMessage id={"unverify"} />
      </DefaultButton>
    );
  }

  return (
    <PrimaryButton
      className={cn(styles.controlButton, styles.verifyButton)}
      theme={verifyButtonTheme}
      onClick={onClickVerify}
    >
      <FormattedMessage id={"verify"} />
    </PrimaryButton>
  );
};

const IdentityListCell: React.FC<IdentityListCellProps> = function IdentityListCell(
  props: IdentityListCellProps
) {
  const {
    identityID,
    oauthProviderType,
    identityType,
    identityName,
    connectedOn,
    addedOn,
    verified,
    toggleVerified,
    onRemoveClicked: _onRemoveClicked,
  } = props;

  const icon = getIcon(identityType, oauthProviderType);

  const onRemoveClicked = useCallback(() => {
    _onRemoveClicked(identityID, identityName ?? "");
  }, [identityID, identityName, _onRemoveClicked]);

  const onVerifyClicked = useCallback(
    (verified: boolean) => {
      toggleVerified?.(identityID, verified);
    },
    [toggleVerified, identityID]
  );

  return (
    <ListCellLayout className={styles.cellContainer}>
      <div className={styles.cellIcon}>{icon}</div>
      <Text className={styles.cellName}>{identityName ?? ""}</Text>
      {verified != null && (
        <>
          {verified ? (
            <Text className={styles.cellDescVerified}>
              <FormattedMessage id="verified" />
            </Text>
          ) : (
            <Text className={styles.cellDescUnverified}>
              <FormattedMessage id="unverified" />
            </Text>
          )}
          <Text className={styles.cellDescSeparator}>{" | "}</Text>
        </>
      )}
      <Text className={styles.cellDesc}>
        {connectedOn != null && (
          <FormattedMessage
            id="UserDetails.connected-identities.connected-on"
            values={{ datetime: connectedOn }}
          />
        )}
        {addedOn != null && (
          <FormattedMessage
            id="UserDetails.connected-identities.added-on"
            values={{ datetime: addedOn }}
          />
        )}
      </Text>
      {verified != null && toggleVerified != null && (
        <VerifyButton verified={verified} toggleVerified={onVerifyClicked} />
      )}
      <DefaultButton
        className={cn(styles.controlButton, styles.removeButton)}
        theme={destructiveTheme}
        onClick={onRemoveClicked}
      >
        <FormattedMessage id={removeButtonTextId[identityType]} />
      </DefaultButton>
    </ListCellLayout>
  );
};

const UserDetailsConnectedIdentities: React.FC<UserDetailsConnectedIdentitiesProps> = function UserDetailsConnectedIdentities(
  props: UserDetailsConnectedIdentitiesProps
) {
  const { identities, availableLoginIdIdentities } = props;
  const { locale, renderToString } = useContext(Context);
  const navigate = useNavigate();
  const {
    deleteIdentity,
    loading: deletingIdentity,
    error: deleteIdentityError,
  } = useDeleteIdentityMutation();

  const [
    isConfirmationDialogVisible,
    setIsConfirmationDialogVisible,
  ] = useState(false);
  const [isErrorDialogVisible, setIsErrorDialogVisible] = useState(false);

  const [confirmationDialogData, setConfirmationDialogData] = useState<
    ConfirmationDialogData
  >({
    identityID: "",
    identityName: "",
  });
  const [errorDialogData, setErrorDialogData] = useState<ErrorDialogData>({
    message: "",
  });

  const identityLists: IdentityLists = useMemo(() => {
    const oauthIdentityList: OAuthIdentityListItem[] = [];
    const emailIdentityList: EmailIdentityListItem[] = [];
    const phoneIdentityList: PhoneIdentityListItem[] = [];
    const usernameIdentityList: UsernameIdentityListItem[] = [];

    // TODO: get actual verified state
    for (const identity of identities) {
      const createdAtStr = formatDatetime(locale, identity.createdAt) ?? "";
      if (identity.type === "OAUTH") {
        const providerType =
          identity.claims["https://authgear.com/claims/oauth/provider_type"];
        oauthIdentityList.push({
          id: identity.id,
          name: identity.claims.email,
          providerType,
          verified: false,
          connectedOn: createdAtStr,
        });
      }

      if (identity.type === "LOGIN_ID") {
        if (
          identity.claims["https://authgear.com/claims/login_id/type"] ===
          "email"
        ) {
          emailIdentityList.push({
            id: identity.id,
            email: identity.claims.email,
            verified: true,
            addedOn: createdAtStr,
          });
        }

        if (
          identity.claims["https://authgear.com/claims/login_id/type"] ===
          "phone"
        ) {
          phoneIdentityList.push({
            id: identity.id,
            phone: identity.claims.phone_number,
            verified: false,
            addedOn: createdAtStr,
          });
        }

        if (
          identity.claims["https://authgear.com/claims/login_id/type"] ===
          "username"
        ) {
          usernameIdentityList.push({
            id: identity.id,
            username: identity.claims.preferred_username,
            addedOn: createdAtStr,
          });
        }
      }
    }
    return {
      oauth: oauthIdentityList,
      email: emailIdentityList,
      phone: phoneIdentityList,
      username: usernameIdentityList,
    };
  }, [locale, identities]);

  const onRemoveClicked = useCallback(
    (identityID: string, identityName: string) => {
      setConfirmationDialogData({
        identityID,
        identityName,
      });
      setIsConfirmationDialogVisible(true);
    },
    [setConfirmationDialogData]
  );

  const onDismissConfirmationDialog = useCallback(() => {
    setIsConfirmationDialogVisible(false);
  }, []);

  const onConfirmRemoveIdentity = useCallback(() => {
    const { identityID } = confirmationDialogData;
    deleteIdentity(identityID).finally(() => {
      onDismissConfirmationDialog();
    });
  }, [confirmationDialogData, deleteIdentity, onDismissConfirmationDialog]);

  useEffect(() => {
    const fallbackErrorMessageId =
      "UserDetails.connected-identities.remove-identity-error.generic";
    const violations = parseError(deleteIdentityError);
    const {
      errorMessageIds,
      unknownViolations,
    } = getErrorMessageIdsFromViolation(violations);

    let errorMessage = null;
    if (errorMessageIds.length > 0) {
      errorMessage = errorMessageIds.map((id) => renderToString(id)).join("\n");
    } else if (unknownViolations.length > 0) {
      errorMessage = renderToString(fallbackErrorMessageId);
    }

    if (errorMessage != null) {
      setErrorDialogData({
        message: errorMessage,
      });
      setIsErrorDialogVisible(true);
    }
  }, [deleteIdentityError, renderToString]);

  const onDismissErrorDialog = useCallback(() => {
    setIsErrorDialogVisible(false);
  }, []);

  const onRenderOAuthIdentityCell = useCallback(
    (item?: OAuthIdentityListItem, _index?: number): React.ReactNode => {
      if (item == null) {
        return null;
      }
      return (
        <IdentityListCell
          identityID={item.id}
          oauthProviderType={item.providerType}
          identityType="oauth"
          identityName={item.name}
          verified={item.verified}
          connectedOn={item.connectedOn}
          onRemoveClicked={onRemoveClicked}
          toggleVerified={() => {}}
        />
      );
    },
    [onRemoveClicked]
  );

  const onRenderEmailIdentityCell = useCallback(
    (item?: EmailIdentityListItem, _index?: number): React.ReactNode => {
      if (item == null) {
        return null;
      }
      return (
        <IdentityListCell
          identityID={item.id}
          identityType="email"
          identityName={item.email}
          verified={item.verified}
          addedOn={item.addedOn}
          onRemoveClicked={onRemoveClicked}
          toggleVerified={() => {}}
        />
      );
    },
    [onRemoveClicked]
  );

  const onRenderPhoneIdentityCell = useCallback(
    (item?: PhoneIdentityListItem, _index?: number): React.ReactNode => {
      if (item == null) {
        return null;
      }
      return (
        <IdentityListCell
          identityID={item.id}
          identityType="phone"
          identityName={item.phone}
          verified={item.verified}
          addedOn={item.addedOn}
          onRemoveClicked={onRemoveClicked}
          toggleVerified={() => {}}
        />
      );
    },
    [onRemoveClicked]
  );

  const onRenderUsernameIdentityCell = useCallback(
    (item?: UsernameIdentityListItem, _index?: number): React.ReactNode => {
      if (item == null) {
        return null;
      }
      return (
        <IdentityListCell
          identityID={item.id}
          identityType="username"
          identityName={item.username}
          addedOn={item.addedOn}
          onRemoveClicked={onRemoveClicked}
        />
      );
    },
    [onRemoveClicked]
  );

  const addIdentitiesMenuProps: IContextualMenuProps = useMemo(() => {
    const availableMenuItem = [
      {
        key: "email",
        text: renderToString("UserDetails.connected-identities.email"),
        iconProps: { iconName: "Mail" },
        onClick: () => navigate("./add-email"),
      },
      {
        key: "phone",
        text: renderToString("UserDetails.connected-identities.phone"),
        iconProps: { iconName: "CellPhone" },
        onClick: () => navigate("./add-phone"),
      },
      {
        key: "username",
        text: renderToString("UserDetails.connected-identities.username"),
        iconProps: { iconName: "Accounts" },
        onClick: () => navigate("./add-username"),
      },
    ];
    const enabledItems = availableMenuItem.filter((item) => {
      return availableLoginIdIdentities.includes(item.key);
    });
    return {
      items: enabledItems,
      directionalHintFixed: true,
    };
  }, [renderToString, navigate, availableLoginIdIdentities]);

  return (
    <div className={styles.root}>
      <Dialog
        hidden={!isConfirmationDialogVisible}
        title={
          <FormattedMessage id="UserDetails.connected-identities.confirm-remove-identity-title" />
        }
        subText={renderToString(
          "UserDetails.connected-identities.confirm-remove-identity-message",
          { identityName: confirmationDialogData.identityName }
        )}
        onDismiss={onDismissConfirmationDialog}
      >
        <DialogFooter>
          <ButtonWithLoading
            labelId="confirm"
            onClick={onConfirmRemoveIdentity}
            loading={deletingIdentity}
          />
        </DialogFooter>
      </Dialog>
      <Dialog
        hidden={!isErrorDialogVisible}
        title={
          <FormattedMessage id="UserDetails.connected-identities.error-dialog-title" />
        }
        subText={errorDialogData.message}
        onDismiss={onDismissErrorDialog}
      >
        <DialogFooter>
          <PrimaryButton onClick={onDismissErrorDialog}>
            <FormattedMessage id="ok" />
          </PrimaryButton>
        </DialogFooter>
      </Dialog>
      <section className={styles.headerSection}>
        <Text as="h2" className={styles.header}>
          <FormattedMessage id="UserDetails.connected-identities.title" />
        </Text>
        <PrimaryButton
          disabled={addIdentitiesMenuProps.items.length === 0}
          iconProps={{ iconName: "CirclePlus" }}
          menuProps={addIdentitiesMenuProps}
          styles={{
            menuIcon: { paddingLeft: "3px" },
            icon: { paddingRight: "3px" },
          }}
        >
          <FormattedMessage id="UserDetails.connected-identities.add-identity" />
        </PrimaryButton>
      </section>
      <section className={styles.identityLists}>
        {identityLists.oauth.length > 0 && (
          <>
            <Text as="h3" className={styles.subHeader}>
              <FormattedMessage id="UserDetails.connected-identities.oauth" />
            </Text>
            <List
              className={styles.list}
              items={identityLists.oauth}
              onRenderCell={onRenderOAuthIdentityCell}
            />
          </>
        )}
        {identityLists.email.length > 0 && (
          <>
            <Text as="h3" className={styles.subHeader}>
              <FormattedMessage id="UserDetails.connected-identities.email" />
            </Text>
            <List
              className={styles.list}
              items={identityLists.email}
              onRenderCell={onRenderEmailIdentityCell}
            />
          </>
        )}
        {identityLists.phone.length > 0 && (
          <>
            <Text as="h3" className={styles.subHeader}>
              <FormattedMessage id="UserDetails.connected-identities.phone" />
            </Text>
            <List
              className={styles.list}
              items={identityLists.phone}
              onRenderCell={onRenderPhoneIdentityCell}
            />
          </>
        )}
        {identityLists.username.length > 0 && (
          <>
            <Text as="h3" className={styles.subHeader}>
              <FormattedMessage id="UserDetails.connected-identities.username" />
            </Text>
            <List
              className={styles.list}
              items={identityLists.username}
              onRenderCell={onRenderUsernameIdentityCell}
            />
          </>
        )}
      </section>
      <Text as="h2" className={styles.primaryIdentitiesTitle}>
        <FormattedMessage id="UserDetails.connected-identities.primary-identities.title" />
      </Text>
      <PrimaryIdentitiesSelectionForm
        className={styles.primaryIdentitiesForm}
        identityLists={identityLists}
      />
    </div>
  );
};

export default UserDetailsConnectedIdentities;
