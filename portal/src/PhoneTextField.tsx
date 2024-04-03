/* global IntlTelInputInitOptions IntlTelInputInstance JSX */
import React, { createRef } from "react";
import { Text, Label } from "@fluentui/react";
import intlTelInput from "intl-tel-input";
import { SystemConfigContext } from "./context/SystemConfigContext";
import styles from "./PhoneTextField.module.css";

export interface PhoneTextFieldValues {
  // Suppose the input now looks like
  // +852 23
  //
  // then
  //
  // rawInputValue === "23"
  // e164 === undefined
  // partialValue === "+85223"
  // alpha2 === "HK"
  // countryCallingCode === "852"
  rawInputValue: string;
  e164?: string;
  partialValue?: string;
  alpha2?: string;
  countryCallingCode?: string;
}

export interface PhoneTextFieldProps {
  className?: string;
  label?: string;
  disabled?: boolean;
  pinnedList?: string[];
  allowlist?: string[];
  initialCountry?: string;
  inputValue: string;
  onChange: (values: PhoneTextFieldValues) => void;
  errorMessage?: React.ReactNode;
}

function makePartialValue(
  countryCallingCode: string,
  rawInputValue: string
): string {
  return `+${countryCallingCode}${rawInputValue}`;
}

export default class PhoneTextField extends React.Component<PhoneTextFieldProps> {
  inputRef: React.RefObject<HTMLInputElement>;
  instance: IntlTelInputInstance | null;

  static contextType = SystemConfigContext;
  // eslint-disable-next-line react/static-property-placement
  declare context: React.ContextType<typeof SystemConfigContext>;

  constructor(props: PhoneTextFieldProps) {
    super(props);
    this.inputRef = createRef();
    this.instance = null;
  }

  componentDidMount(): void {
    const options: IntlTelInputInitOptions = {
      autoPlaceholder: "aggressive",
      customContainer: styles.container,
    };
    if (this.props.initialCountry != null) {
      options.initialCountry = this.props.initialCountry;
    }
    if (this.props.allowlist != null) {
      options.onlyCountries = [...this.props.allowlist];
    }
    if (this.props.pinnedList != null) {
      options.preferredCountries = [...this.props.pinnedList];
    } else {
      options.preferredCountries = [];
    }

    if (this.inputRef.current != null) {
      const instance = intlTelInput(this.inputRef.current, options);
      instance.setNumber(this.props.inputValue);
      this.instance = instance;

      this.inputRef.current.addEventListener("input", this.onInputChange);
      this.inputRef.current.addEventListener(
        "countrychange",
        this.onCountryChange
      );
    }
  }

  componentDidUpdate(prevProps: PhoneTextFieldProps): void {
    if (prevProps.inputValue !== this.props.inputValue) {
      this.instance?.setNumber(this.props.inputValue);
    }
  }

  onInputChange = (): void => {
    this.emitOnChange();
  };

  onCountryChange = (): void => {
    this.emitOnChange();
  };

  emitOnChange(): void {
    if (this.instance != null && this.inputRef.current != null) {
      const rawInputValue = this.inputRef.current.value;
      const countryData = this.instance.getSelectedCountryData();
      const alpha2 = countryData.iso2;
      const countryCallingCode = countryData.dialCode;
      // The output of getNumber() is very unstable.
      // If isPossibleNumber(), then it has +countryCallingCode,
      // otherwise it is rawInputValue with some spaces in it.
      const maybeInvalid = this.instance.getNumber();
      let e164;
      if (this.instance.isPossibleNumber()) {
        if (maybeInvalid != null) {
          e164 = maybeInvalid;
        }
      }
      let partialValue;
      if (e164 != null) {
        partialValue = e164;
      } else if (countryCallingCode != null) {
        partialValue = makePartialValue(countryCallingCode, rawInputValue);
      }

      const values = {
        e164,
        rawInputValue,
        alpha2,
        countryCallingCode,
        partialValue,
      };

      this.props.onChange(values);
    }
  }

  render(): JSX.Element {
    const { className, label, errorMessage, disabled } = this.props;
    const semanticColors = this.context?.themes.main.semanticColors;
    const inputBorder = semanticColors?.inputBorder ?? "";
    const errorText = semanticColors?.errorText ?? "";
    const inputFocusBorderAlt = semanticColors?.inputFocusBorderAlt ?? "";
    const disabledBackground = semanticColors?.disabledBackground ?? "";
    return (
      <div className={className}>
        {label ? <Label disabled={disabled}>{label}</Label> : null}
        <input
          style={{
            // @ts-expect-error
            "--PhoneTextField-border-color":
              errorMessage != null ? errorText : inputBorder,
            "--PhoneTextField-border-color-focus":
              errorMessage != null ? errorText : inputFocusBorderAlt,
            backgroundColor: disabled ? disabledBackground : undefined,
          }}
          className={styles.input}
          type="text"
          ref={this.inputRef}
          disabled={disabled}
        />
        {errorMessage ? (
          <Text
            block={true}
            styles={{
              root: {
                color: errorText,
              },
            }}
            className={styles.errorMessage}
          >
            {errorMessage}
          </Text>
        ) : null}
      </div>
    );
  }
}
