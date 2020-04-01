package goex

func ModifyItem(input ModTemplate, realizedEdit bool, affectTemplate bool) error {
	item, err := GetSpecificActiveItem(input.ID)

	if err != nil {
		return err
	}

	template, err := GetSpecificTemplate(input.ID)

	if err != nil {
		return err
	}

	// update each trait if non-default value were passed
	if input.Amount != 0 {
		template.Amount = input.Amount
	}
	if input.Category != "" {
		template.Category = input.Category
	}
	if input.Description != "" {
		template.Description = input.Description
	}
	if input.Name != "" {
		template.Name = input.Name
	}
	if realizedEdit {
		item.Realized = input.Realized
	}

	if affectTemplate {
		_ = DeleteTemplateItem(input.ID)
	}
	_ = DeleteActiveItem(input.ID)

	if affectTemplate {
    	_, _ = WriteNewTemplate(&template, false)
	}
    _ = WriteNewMonthItem(&template, item.Realized)

    return nil
}