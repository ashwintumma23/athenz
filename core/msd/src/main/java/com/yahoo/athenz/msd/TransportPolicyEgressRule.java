//
// This file generated by rdl 1.5.2. Do not modify!
//

package com.yahoo.athenz.msd;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.yahoo.rdl.*;

//
// TransportPolicyEgressRule - Transport policy egress rule
//
@JsonIgnoreProperties(ignoreUnknown = true)
public class TransportPolicyEgressRule {
    public long id;
    public Timestamp lastModified;
    public TransportPolicyEntitySelector entitySelector;
    public TransportPolicyPeer to;

    public TransportPolicyEgressRule setId(long id) {
        this.id = id;
        return this;
    }
    public long getId() {
        return id;
    }
    public TransportPolicyEgressRule setLastModified(Timestamp lastModified) {
        this.lastModified = lastModified;
        return this;
    }
    public Timestamp getLastModified() {
        return lastModified;
    }
    public TransportPolicyEgressRule setEntitySelector(TransportPolicyEntitySelector entitySelector) {
        this.entitySelector = entitySelector;
        return this;
    }
    public TransportPolicyEntitySelector getEntitySelector() {
        return entitySelector;
    }
    public TransportPolicyEgressRule setTo(TransportPolicyPeer to) {
        this.to = to;
        return this;
    }
    public TransportPolicyPeer getTo() {
        return to;
    }

    @Override
    public boolean equals(Object another) {
        if (this != another) {
            if (another == null || another.getClass() != TransportPolicyEgressRule.class) {
                return false;
            }
            TransportPolicyEgressRule a = (TransportPolicyEgressRule) another;
            if (id != a.id) {
                return false;
            }
            if (lastModified == null ? a.lastModified != null : !lastModified.equals(a.lastModified)) {
                return false;
            }
            if (entitySelector == null ? a.entitySelector != null : !entitySelector.equals(a.entitySelector)) {
                return false;
            }
            if (to == null ? a.to != null : !to.equals(a.to)) {
                return false;
            }
        }
        return true;
    }
}
