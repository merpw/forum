import React from "react";
import { useState } from "react";
import Select from "react-select";

export default function CategoryMultiselect() {
  const [category, setCategory] = useState([]);
  
  const optionsCategory = [
    { value: "facts", label: "Facts" },
    { value: "rumors", label: "Rumors" }
  ];

  const handleCategoryChange = async (selected, selectaction) => {
    const { action } = selectaction;
    // console.log(`action ${action}`);
    if (action === "clear") {
    } else if (action === "select-option") {
    } else if (action === "remove-value") {
      console.log("remove");
    }
    setCategory(selected);
  };

  return (
    <div>
      <Select
        id="selectCategory"
        instanceId="selectCategory"
        isMulti
        name="categories"
        className="categoryMultiselect"
        classNamePrefix="select"
        options={optionsCategory}
        onChange={handleCategoryChange}
        placeholder="Category"
      />
      <br />
      <br />
    </div>
  );
}
