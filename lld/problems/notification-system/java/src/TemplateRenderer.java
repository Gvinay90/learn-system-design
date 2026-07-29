import java.util.Map;

/** Basic "{placeholder}" template substitution. */
public final class TemplateRenderer {
    private TemplateRenderer() {}

    /**
     * Substitutes "{key}" tokens in template with data.get(key). Tokens with
     * no matching key are left untouched.
     */
    public static String render(String template, Map<String, String> data) {
        StringBuilder out = new StringBuilder();
        int i = 0;
        while (i < template.length()) {
            int open = template.indexOf('{', i);
            if (open == -1) {
                out.append(template.substring(i));
                break;
            }
            int close = template.indexOf('}', open);
            if (close == -1) {
                out.append(template.substring(i));
                break;
            }
            out.append(template, i, open);
            String key = template.substring(open + 1, close);
            if (data != null && data.containsKey(key)) {
                out.append(data.get(key));
            } else {
                out.append(template, open, close + 1);
            }
            i = close + 1;
        }
        return out.toString();
    }
}
