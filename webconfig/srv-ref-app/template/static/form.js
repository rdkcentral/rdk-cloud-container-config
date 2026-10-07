document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('webconfig-form');
    const feedbackDiv = document.getElementById('feedback');
    const submitBtn = document.getElementById('submit-btn');

    /**
     * Validate JSON string
     */
    function validateJSON(jsonString) {
        try {
            JSON.parse(jsonString);
            return true;
        } catch (e) {
            return false;
        }
    }

    /**
     * Show feedback message
     */
    function showFeedback(message, type = 'success') {
        feedbackDiv.textContent = message;
        feedbackDiv.className = `feedback ${type}`;
        feedbackDiv.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }

    /**
     * Hide feedback message
     */
    function hideFeedback() {
        feedbackDiv.className = 'feedback';
        feedbackDiv.textContent = '';
    }

    /**
     * Validate form inputs
     */
    function validateForm() {
        const subdocName = document.getElementById('subdoc_name').value.trim();
        const subdocData = document.getElementById('subdoc_data').value.trim();
        const paramName = document.getElementById('param_name').value.trim();
        const macAddress = document.getElementById('mac_address').value.trim();

        // Check for empty fields
        if (!subdocName || !subdocData || !paramName || !macAddress) {
            showFeedback('All fields are required', 'error');
            return false;
        }
        // Validate JSON
        if (!validateJSON(subdocData)) {
            showFeedback('Configuration Data must be valid JSON', 'error');
            return false;
        }

        // Validate subdoc name format
        if (!/^[a-zA-Z0-9_-]+$/.test(subdocName)) {
            showFeedback('Subdocument name contains invalid characters', 'error');
            return false;
        }

        // Validate TR181 parameter format
        if (!/^Device\.[a-zA-Z0-9_.]+$/.test(paramName)) {
            showFeedback('TR181 parameter must follow format: Device.Subsystem.Parameter', 'error');
            return false;
        }

        // Validate MAC address format
        if (!/^([0-9a-fA-F]{2}:){5}([0-9a-fA-F]{2})$|^([0-9a-fA-F]{2}){6}$/.test(macAddress)) {
            showFeedback('MAC address must be in format XX:XX:XX:XX:XX:XX or XXXXXXXXXXXX', 'error');
            return false;
        }

        return true;
    }

    /**
     * Handle form submission
     */
    form.addEventListener('submit', async function(e) {
        e.preventDefault();

        // Clear previous feedback
        hideFeedback();

        // Validate form
        if (!validateForm()) {
            return;
        }

        // Disable submit button and show loading state
        submitBtn.disabled = true;
        const originalText = submitBtn.textContent;
        submitBtn.innerHTML = '<span class="spinner"></span>Submitting...';

        try {
            // Prepare form data
            const formData = new FormData(form);

            // Send request
            const response = await fetch('/app1/send', {
                method: 'POST',
                body: formData
            });

            // Parse response
            let data;
            try {
                data = await response.json();
            } catch (e) {
                data = { message: await response.text() };
            }

            // Handle response
            if (response.ok) {
                showFeedback(
                    `✓ ${data.message || 'Configuration submitted successfully'}`,
                    'success'
                );

                // Reset form after success
                setTimeout(() => {
                    form.reset();
                    hideFeedback();
                }, 3000);
            } else {
                showFeedback(
                    `✗ ${data.error || data.message || 'Failed to submit configuration'}`,
                    'error'
                );
            }
        } catch (error) {
            showFeedback(
                `✗ Network error: ${error.message}`,
                'error'
            );
        } finally {
            // Re-enable submit button
            submitBtn.disabled = false;
            submitBtn.textContent = originalText;
        }
    });

    // Real-time JSON validation
    const subdocDataInput = document.getElementById('subdoc_data');
    subdocDataInput.addEventListener('change', function() {
        if (this.value.trim() && !validateJSON(this.value)) {
            this.setCustomValidity('Invalid JSON');
        } else {
            this.setCustomValidity('');
        }
    });
});
