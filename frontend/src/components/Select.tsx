import ReactSelect from "react-select"

const Select: ReactSelect = (props) => {
  return (
    <ReactSelect
      {...props}
      unstyled={true}
      className={""}
      styles={{
        control: (provided) => ({
          ...provided,
          transition: undefined,
          outline: undefined,
        }),
      }}
      classNames={{
        control: ({ isFocused }) =>
          "relative input " +
          " " +
          (isFocused ? "outline outline-offset-2 outline-base-content/20" : "text-gray-400"),

        placeholder: () => "text-sm text-gray-400",

        menu: () => "bg-base-300 my-1 mt-1 rounded",
        option: ({ isFocused }) => "p-2" + " " + (isFocused ? "bg-base-200" : "bg-base-100"),
        noOptionsMessage: () => "my-3",

        indicatorsContainer: () => "absolute w-full flex justify-end h-11 pr-4",
        indicatorSeparator: () => "my-2 bg-gray-400",

        dropdownIndicator: () => "p-2",
        clearIndicator: () => "p-2",

        valueContainer: () => "absolute flex-none gap-1 overflow-scroll",
        multiValue: () => "bg-base-300 text-base-content rounded-md my-auto overflow-items-scroll",
        multiValueLabel: () => "ml-2 my-1",
        multiValueRemove: () => "hover:bg-error px-2 rounded-md",
      }}
    />
  )
}

export default Select
