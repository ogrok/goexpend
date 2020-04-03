package state

import (
	"github.com/adaminoue/goexpend/src/models"
)

func ModifyItem(input *models.Modification, realizedEdit bool, affectTemplate bool) error {
	var template models.Template

	// it is valid if template does not exist; one-time items do not have them
	template, _ = GetSpecificTemplate(input.ID)
	item, err := GetSpecificActiveItem(input.ID)

	if err != nil {
		return err
	}

	// update each trait if non-default value was passed
	if input.Amount != 0 {
		template.Amount = input.Amount
		item.Amount = input.Amount
	}
	if input.Category != "" {
		template.Category = input.Category
		item.Category = input.Category
	}
	if input.Description != "" {
		template.Description = input.Description
		item.Description = input.Description
	}
	if input.Name != "" {
		template.Name = input.Name
		item.Name = input.Name
	}
	if realizedEdit {
		item.Realized = input.Realized
	}

	// this entire section below is so unsafe and should not exist in this form wow
	if affectTemplate {
		_ = DeleteTemplateItem(input.ID, false)
	}
	_ = DeleteActiveItem(input.ID)

	if affectTemplate {
    	_, _ = WriteNewTemplate(&template, false)
	}
    _ = WriteNewActiveItemDirectly(&item)

    return nil
}