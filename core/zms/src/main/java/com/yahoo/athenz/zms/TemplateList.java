//
// This file generated by rdl 1.4.14. Do not modify!
//

package com.yahoo.athenz.zms;
import java.util.List;
import com.yahoo.rdl.*;

//
// TemplateList - List of template names that is the base struct for server and
// domain templates
//
public class TemplateList {
    public List<String> templateNames;

    public TemplateList setTemplateNames(List<String> templateNames) {
        this.templateNames = templateNames;
        return this;
    }
    public List<String> getTemplateNames() {
        return templateNames;
    }

    @Override
    public boolean equals(Object another) {
        if (this != another) {
            if (another == null || another.getClass() != TemplateList.class) {
                return false;
            }
            TemplateList a = (TemplateList) another;
            if (templateNames == null ? a.templateNames != null : !templateNames.equals(a.templateNames)) {
                return false;
            }
        }
        return true;
    }
}
